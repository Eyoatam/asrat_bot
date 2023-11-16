FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/bot

EXPOSE 8443

CMD ["./bin/bot"]
