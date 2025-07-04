FROM golang:1.24.4-alpine

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o app .

RUN chmod +x entrypoint.sh

EXPOSE 8080

CMD ["./entrypoint.sh"]
