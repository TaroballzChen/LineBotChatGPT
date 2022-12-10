FROM golang:1.19.3

RUN apt update
RUN apt install -y curl git
RUN git config --global http.sslVerify false
RUN cd / && git clone https://github.com/TaroballzChen/LineBotChatGPT

WORKDIR /LineBotChatGPT

EXPOSE 80

ENTRYPOINT ["go","run","main.go"]