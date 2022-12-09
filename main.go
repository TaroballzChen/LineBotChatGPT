package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func callAI(question string) string {
	out, err := exec.Command("python3", "ai.py", "--question", question).Output()
	if err != nil {
		log.Println(err)
	}
	return string(out)
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
				answer := callAI(question)
				answer = strings.Replace(answer, "AI:", "", 1)
				log.Println("A:", answer)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(answer)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
