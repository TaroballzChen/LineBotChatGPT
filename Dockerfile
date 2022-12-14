FROM golang:1.19.3

RUN apt update
RUN apt install -y curl git
RUN git config --global http.sslVerify false
RUN cd / && git clone https://github.com/TaroballzChen/LineBotChatGPT

WORKDIR /LineBotChatGPT

ADD "https://www.random.org/cgi-bin/randbyte?nbytes=10&format=h" skipcache
RUN git pull

EXPOSE 80

ENTRYPOINT ["go","run","main.go"]