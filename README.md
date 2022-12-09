# LineBotChatGPT
utilized openai api for implementation of chatGPT with LineBot

# Work Environment for myself
- `Python 3.10.4`
- `go version go1.18.1 darwin/amd64`

# Usage

```shell
git clone https://github.com/TaroballzChen/LineBotChatGPT
pip3 install openai python-dotenv
cd LineBotChatGPT
touch .env
echo ChannelSecret=your_LINE_ChannelSecret >> .env
echo ChannelAccessToken=your_LINE_ChannelAccessToken >> .env
echo OpenApiKey=your_OpenApiKey >>.env
go run main.go
```

then use `ngrok` or other method(cloudcontainer, nginx with certbot etc.) exposed `80` port to public network with SSL 

enjoy!


# Reference
1. https://github.com/kkdai/linebot-group
2. https://github.com/kkdai/LineBotTemplate