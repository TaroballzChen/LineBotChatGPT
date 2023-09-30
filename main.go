package main

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// linebot client ptr
var bot *linebot.Client

// OpenAI Api key
var OpenAIApiKey string
var AIName string

// CompletionModelParam
var MaxTokens int
var Temperature float32
var TopP float32
var FrequencyPenalty float32
var PresencePenalty float32
var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

// chatWithAI
var isChatWithAnotherAI bool
var chatPartner string

func main() {
	var err error
	OpenAIApiKey = os.Getenv("OpenApiKey")
	AIName = os.Getenv("AIName")
	GetModelParamFromEnv()
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := "80"
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func GetModelParamFromEnv() {
	var err error
	if MaxTokens, err = getenvInt("max_tokens"); err != nil {
		log.Println("max_tokens", err)
		err = nil
	}
	if Temperature, err = getenvFloat("temperature"); err != nil {
		log.Println("temperature", err)
		err = nil
	}
	if TopP, err = getenvFloat("top_p"); err != nil {
		log.Println("top_p", err)
		err = nil
	}
	if FrequencyPenalty, err = getenvFloat("FrequencyPenalty"); err != nil {
		log.Println("FrequencyPenalty", err)
		err = nil
	}
	if PresencePenalty, err = getenvFloat("PresencePenalty"); err != nil {
		log.Println("PresencePenalty", err)
		err = nil
	}
	if isChatWithAnotherAI, err = getenvBoolean("isChatWithAnotherAI"); err != nil {
		log.Println("isChatWithAnotherAI", err)
		err = nil
	}

	if isChatWithAnotherAI {
		if chatPartner, err = getenvStr("chatPartner"); err != nil {
			log.Println("chatPartner", err)
			err = nil
		}
	}
}

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	log.Println(key, v)
	return v, nil
}

func getenvInt(key string) (int, error) {
	s, err := getenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func getenvFloat(key string) (float32, error) {
	s, err := getenvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

func getenvBoolean(key string) (bool, error) {
	s, err := getenvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}

func isChatPartner(user linebot.UserProfileResponse) bool {
	if user.DisplayName == chatPartner {
		log.Println("the chatPartner", chatPartner, "is existed!")
		return true
	}
	return false
}

func GetResponse(client *openai.Client, ctx context.Context, quesiton string) string {
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一個知識儲量非常豐富且有問必答的強大AI助理",
			},
			{
				Role:    "user",
				Content: quesiton,
			},
		},
		MaxTokens:        MaxTokens,
		Temperature:      Temperature,
		TopP:             TopP,
		FrequencyPenalty: FrequencyPenalty,
		PresencePenalty:  PresencePenalty,
	})

	if err != nil {
		log.Println("Get Open AI Response Error: ", err)
	}
	answer := resp.Choices[0].Message.Content
	answer = strings.TrimSpace(answer)
	return answer
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:

				question := message.Text

				switch {
				case event.Source.GroupID != "":
					//In the group
					if !strings.HasPrefix(message.Text, AIName) {
						log.Println("Group", event.Source.GroupID, "message: ", message.Text)
						return
					}
					question = strings.Replace(question, AIName, "", 1)

				case event.Source.RoomID != "":
					//In the room
					if !strings.HasPrefix(message.Text, AIName) {
						log.Println("Room", event.Source.RoomID, "message: ", message.Text)
						return
					}
					question = strings.Replace(question, AIName, "", 1)

				}

				log.Println("Q:", question)

				ctx := context.Background()
				client := openai.NewClient(OpenAIApiKey)
				answer := GetResponse(client, ctx, question)
				log.Println("A:", answer)

				if event.Source.GroupID != "" && isChatWithAnotherAI && chatPartner != "" {
					if profile, err := bot.GetGroupMemberProfile(event.Source.GroupID, event.Source.UserID).Do(); err == nil {
						if isChatPartner(*profile) {
							answer = chatPartner + answer
						}
					}

				}

				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(answer)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
