FROM golang:1.17

ADD ./ /go/src/gdrive-telegram-bot

WORKDIR /go/src/gdrive-telegram-bot/src

RUN env GOOS=linux GOARCH=amd64 go build -o "gdrive-telegram-bot"

CMD ["./gdrive-telegram-bot"]