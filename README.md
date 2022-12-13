# LineBotChatGPT
utilized openai api for implementation of chatGPT with LineBot

## Work Environment for myself
- `go version go1.18.1 darwin/amd64`

## Usage

```shell
git clone https://github.com/TaroballzChen/LineBotChatGPT
cd LineBotChatGPT
echo ChannelSecret=your_LINE_ChannelSecret >> .env
echo ChannelAccessToken=your_LINE_ChannelAccessToken >> .env
echo OpenApiKey=your_OpenApiKey >>.env
# you could modify the GPT-3 completion model parameter by modifying `.env` file
go run main.go
```

then use `ngrok` or other method(cloud container, nginx with certbot etc.) exposed `80` port to public network with SSL 

enjoy!

### Dockerfile
1. download the `Dockerfile`, `.env` file in this project
2. `docker build --no-cache -t LineBotChatGPT:latest .`
3. modify the `.env` file to fill the LINEBOT and OpneAI token:

```shell
echo ChannelSecret=your_LINE_ChannelSecret >> .env
echo ChannelAccessToken=your_LINE_ChannelAccessToken >> .env
echo OpenApiKey=your_OpenApiKey >>.env
# try to modify the GPT-3 completion model parameter by modifying `.env` file
```

4. `docker run -p 8080:80 -v $PWD/.env:/LineBotChatGPT/.env LineBotChatGPT`
then use `ngrok` or other method(cloud container like railway.app, Heroku or nginx with certbot etc.) exposed `80` port to public network with SSL

enjoy!

### Railway.app Deploy
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template/UTbtDt?referralCode=Taroballz)

- fill your LINE `ChannelSecret`, `ChannelAccessToken` and OpenAI `OpenApiKey` token

P.S. You should own the github account to sign up the railway.app account. When you create container by my template on the above, the railway.app would help you fork my github project to your repo. then you could modify the model parameter's value on your forked project.

enjoy!

TODO: [tutorail video]()


## result
![img.png](img.png)

## Reference
1. https://github.com/kkdai/linebot-group
2. https://github.com/kkdai/LineBotTemplate
3. https://www.learncodewithmike.com/2020/06/python-line-bot.html
4. https://github.com/kkdai/chatgpt
5. https://github.com/PullRequestInc/go-gpt3
