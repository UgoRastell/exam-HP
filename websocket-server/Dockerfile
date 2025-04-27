FROM golang:1.24-alpine

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o websocket-server .

EXPOSE 8080
CMD ["./websocket-server"]
