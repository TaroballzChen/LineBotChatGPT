FROM golang:1.19.3

RUN apt update
RUN apt install -y curl git python3-pip
RUN git config --global http.sslVerify false

RUN pip3 install python-dotenv openai

RUN cd / && git clone https://github.com/TaroballzChen/LineBotChatGPT

WORKDIR /LineBotChatGPT

EXPOSE 80

ENTRYPOINT ["go","run","main.go"]