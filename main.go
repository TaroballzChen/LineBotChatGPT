package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/3JoB/anthropic-sdk-go/v2/pkg/pool"
	"io"
	//"github.com/3JoB/anthropic-sdk-go/v2"
	//"github.com/3JoB/anthropic-sdk-go/v2/data"
	//"github.com/3JoB/anthropic-sdk-go/v2/pkg/pool"
	//"github.com/3JoB/anthropic-sdk-go/v2/resp"
	"github.com/bincooo/claude-api"
	"github.com/bincooo/claude-api/types"
	"github.com/google/generative-ai-go/genai"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Claude2Chat object map for different user/group/room
var Claude2Chat = map[string]types.Chat{}

//var Claude2Chat = map[string]*pool.Session{}

// GeminiChat object map for different user/group/room
var GeminiChat = map[string]map[string]interface{}{}

// linebot client ptr
var bot *linebot.Client

// OpenAI Api key
var OpenAIApiKey string
var GPTName string

// Claude2 Api key
var Claude2ApiKey string
var Claude2Name string

// Gemini AI Api key
var GeminiApiKey string
var GeminiName string

// DeepSeek Azure Api key
var DeepSeekAzureApiUrl string
var DeepSeekAzureApiKey string
var DeepSeekName string

// CompletionModelParam
var MaxTokens int
var Temperature float32
var TopP float32
var FrequencyPenalty float32
var PresencePenalty float32
var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func main() {
	var err error
	GeminiApiKey = os.Getenv("GeminiApiKey")
	OpenAIApiKey = os.Getenv("OpenAIApiKey")
	Claude2ApiKey = os.Getenv("Claude2ApiKey")
	DeepSeekAzureApiKey = os.Getenv("DeepSeekAzureApiKey")
	GeminiName = os.Getenv("GeminiName")
	GPTName = os.Getenv("GPTName")
	Claude2Name = os.Getenv("Claude2Name")
	DeepSeekName = os.Getenv("DeepSeekName")
	DeepSeekAzureApiUrl = os.Getenv("DeepSeekAzureApiUrl")
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

func GetResponse(client *openai.Client, ctx context.Context, question string) string {
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一個知識儲量非常豐富、專業且有問必答的強大AI助理",
			},
			{
				Role:    "user",
				Content: question,
			},
		},
		MaxTokens:        MaxTokens,
		Temperature:      Temperature,
		TopP:             TopP,
		FrequencyPenalty: FrequencyPenalty,
		PresencePenalty:  PresencePenalty,
	})

	if err != nil {
		errString := fmt.Sprintf("Get Open AI Response Error: %s", err)
		log.Println(errString)
		return errString
	}
	answer := resp.Choices[0].Message.Content
	answer = strings.TrimSpace(answer)
	return answer
}

func GetImageResponse(client *openai.Client, ctx context.Context, question string) string {
	resp, err := client.CreateImage(ctx, openai.ImageRequest{
		Prompt:         question,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	},
	)

	if err != nil {
		errString := fmt.Sprintf("Image creation error: %s", err)
		log.Println("Image creation error: ", err)
		return errString
	}

	answer := resp.Data[0].URL
	return answer
}

func Cluaude2OutputText(partialResponse chan types.PartialResponse) string {
	var text string
	for {
		message, ok := <-partialResponse
		if !ok {
			return text
		}

		if message.Error != nil {
			log.Println("Claude2 join message Error:", message.Error)
		}

		text += message.Text
	}
}

func Cluaude3OutputText(Session *pool.Session) string {
	return Session.Response.String()
}

func GeminiOutputText(ChatSession *genai.ChatSession, ctx context.Context, question string) string {
	res, err := ChatSession.SendMessage(ctx, genai.Text(question))
	if err != nil {
		errString := fmt.Sprintf("Gemini Response Error: %s", err)
		log.Println(errString)
		return errString
	}

	var text string
	for _, cand := range res.Candidates {
		for _, part := range cand.Content.Parts {
			part_string := fmt.Sprintf("%s", part)
			text += part_string
		}
	}
	return text
}

func DeepSeekOutputText(question string) string {
	data := fmt.Sprintf(`{"model": "DeepSeek-R1","messages": [{"role": "user","content": "%s"}]}`, question)
	jsonStr := []byte(data)
	req, err := http.NewRequest("POST", DeepSeekAzureApiUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", DeepSeekAzureApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("DeepSeek Response Error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("DeepSeek Response Error:", resp.Status)
		return "DeepSeek Response Error:" + resp.Status
	}

	type Response struct {
		Choices []struct {
			FinishReason string `json:"finish_reason"`
			Index        int    `json:"index"`
			Message      struct {
				Content          string `json:"content"`
				ReasoningContent any    `json:"reasoning_content"`
				Role             string `json:"role"`
				ToolCalls        any    `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
		Created             int    `json:"created"`
		ID                  string `json:"id"`
		Model               string `json:"model"`
		Object              string `json:"object"`
		UsagempletionTokens int    `json:"usagempletion_tokens"`
		PromptTokens        int    `json:"prompt_tokens"`
		PromptTokensDetails any    `json:"prompt_tokens_details"`
		TotalTokens         int    `json:"total_tokens"`
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	var serializeResp Response
	if err := json.Unmarshal(body, &serializeResp); err != nil { // Parse []byte to go struct pointer
		fmt.Println("DeepSeek Response Error: Can not unmarshal JSON")
	}
	text := serializeResp.Choices[0].Message.Content
	text = strings.Replace(text, "<think>", "", 1)
	text = strings.Replace(text, "</think>", "", 1)
	text = strings.TrimSpace(text)
	return text
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
				var _ID string
				switch {
				case event.Source.GroupID != "":
					//In the group
					_ID = event.Source.GroupID
					if !strings.HasPrefix(message.Text, GPTName) && !strings.HasPrefix(message.Text, Claude2Name) && !strings.HasPrefix(message.Text, GeminiName) && !strings.HasPrefix(message.Text, DeepSeekName) {
						log.Println("Group", event.Source.GroupID, "message: ", message.Text)
						return
					}

				case event.Source.RoomID != "":
					//In the room
					_ID = event.Source.RoomID
					if !strings.HasPrefix(message.Text, GPTName) && !strings.HasPrefix(message.Text, Claude2Name) && !strings.HasPrefix(message.Text, GeminiName) && !strings.HasPrefix(message.Text, DeepSeekName) {
						log.Println("Room", event.Source.RoomID, "message: ", message.Text)
						return
					}
				case event.Source.UserID != "":
					//In the personal chat
					_ID = event.Source.UserID
				}

				// decide the AI object
				var AIObject string
				switch {
				case strings.HasPrefix(message.Text, GPTName):
					AIObject = GPTName
				case strings.HasPrefix(message.Text, Claude2Name):
					AIObject = Claude2Name
				case strings.HasPrefix(message.Text, GeminiName):
					AIObject = GeminiName
				case strings.HasPrefix(message.Text, DeepSeekName):
					AIObject = DeepSeekName
				default:
					AIObject = DeepSeekName
				}
				question = strings.Replace(question, AIObject, "", 1)

				log.Println("Q:", question)

				ctx := context.Background()

				var answer string

				switch AIObject {
				case DeepSeekName:
					answer = DeepSeekOutputText(question)
				case GPTName:
					client := openai.NewClient(OpenAIApiKey)
					// Handle the special question with image and text
					switch {
					case strings.HasPrefix(question, "作圖"):
						question = strings.Replace(question, "作圖", "", 1)
						answer = GetImageResponse(client, ctx, question)
					default:
						answer = GetResponse(client, ctx, question)
					}
				case Claude2Name:
					//options := anthropic.Config{Key: Claude2ApiKey, DefaultModel: data.ModelFullClaude3Haiku}
					//chat, err := anthropic.New(&options)
					//if err != nil {
					//	log.Println("New Claude Chat Error:", err)
					//}
					if _, ok := Claude2Chat[_ID]; !ok {
						//old package usage
						options := claude.NewDefaultOptions(Claude2ApiKey, "", "claude-3-sonnet-20240229")
						chatObj, err := claude.New(options)

						//d, err := chat.Send(
						//	&anthropic.Sender{
						//		Message: data.MessageModule{
						//			Human: "你現在是一個強大且知識量儲備非常豐富的AI助理",
						//		},
						//		Sender: &resp.Sender{MaxToken: uint(MaxTokens)},
						//	})
						if err != nil {
							log.Println("New Claude2 Chat Error:", err)
						}

						Claude2Chat[_ID] = chatObj
					}
					if question == "銷毀記憶" {
						Claude2Chat[_ID].Delete()
						//chat.CloseSession(Claude2Chat[_ID])
						delete(Claude2Chat, _ID)
						answer = "已銷毀編號為" + _ID + "的記憶"
					} else {
						chat := Claude2Chat[_ID]
						partialResponse, err := chat.Reply(ctx, question, nil)

						//ds, err := chat.Send(
						//	&anthropic.Sender{
						//		Message: data.MessageModule{
						//			Human: question,
						//		},
						//		SessionID: Claude2Chat[_ID].ID,
						//		Sender:    &resp.Sender{MaxToken: uint(MaxTokens)},
						//	})
						if err != nil {
							log.Println("Call Claude API and occur response error:", err)
						}
						answer = Cluaude2OutputText(partialResponse)
						//answer = Cluaude3OutputText(ds)
					}
				case GeminiName:
					if _, ok := GeminiChat[_ID]; !ok {
						client, err := genai.NewClient(ctx, option.WithAPIKey(GeminiApiKey))
						if err != nil {
							log.Println("New Gemini Chat Error:", err)
						}
						GeminiChat[_ID] = map[string]interface{}{}
						GeminiChat[_ID]["client"] = client
						model := client.GenerativeModel("gemini-pro")
						cs := model.StartChat()
						GeminiChat[_ID]["chatSession"] = cs
					}
					if question == "銷毀記憶" {
						err := GeminiChat[_ID]["client"].(*genai.Client).Close()
						if err != nil {
							log.Println("Close Gemini Chat Error:", err)
						}
						delete(GeminiChat, _ID)
						answer = "已銷毀編號為" + _ID + "的Gemini記憶"
					} else {
						cs := GeminiChat[_ID]["chatSession"].(*genai.ChatSession)
						answer = GeminiOutputText(cs, ctx, question)
					}
				}

				log.Println("A:", answer)

				switch {
				case strings.HasPrefix(answer, "https://") && AIObject == GPTName:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(answer, answer)).Do(); err != nil {
						log.Print(err)
					}
				default:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(answer)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
