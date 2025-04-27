FROM golang:1.24-alpine

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o pokemon-client-grpc .

CMD ["./pokemon-client-grpc"]
