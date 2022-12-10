package main

import (
	"context"
	"fmt"
	gpt3 "github.com/PullRequestInc/go-gpt3"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
	"os"
	"strings"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := "80"
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func GetResponse(client gpt3.Client, ctx context.Context, quesiton string) string {
	resp, err := client.CompletionWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:        gpt3.IntPtr(3000),
		Temperature:      gpt3.Float32Ptr(0.9),
		TopP:             gpt3.Float32Ptr(1),
		FrequencyPenalty: float32(0),
		PresencePenalty:  float32(0.6),
	})
	if err != nil {
		log.Println("Get Open AI Response Error: ", err)
	}
	answer := resp.Choices[0].Text
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
				AIName := "小A"

				if !strings.HasPrefix(message.Text, AIName) {
					log.Println("message: ", message.Text)
					return
				}

				question := strings.Replace(message.Text, AIName, "", 1)
				log.Println("Q:", question)

				apiKey := os.Getenv("OpenApiKey")
				if apiKey == "" {
					panic("Missing API KEY")
				}
				ctx := context.Background()
				client := gpt3.NewClient(apiKey)
				answer := GetResponse(client, ctx, question)
				answer = strings.Replace(answer, "AI:", "", 1)
				log.Println("A:", answer)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(answer)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
