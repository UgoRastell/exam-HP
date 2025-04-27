FROM golang:1.24-alpine

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o pokemon-client-websocket .

EXPOSE 8080
CMD ["./pokemon-client-websocket"]
