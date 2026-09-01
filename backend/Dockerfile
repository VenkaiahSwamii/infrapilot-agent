FROM golang:1.26.4

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o infrapilot ./cmd/server

EXPOSE 8080

CMD ["./infrapilot"]
