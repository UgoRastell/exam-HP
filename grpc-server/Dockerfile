FROM golang:1.24-alpine

WORKDIR /app

COPY . .
RUN go mod tidy
RUN go build -o grpc-server .

EXPOSE 50051
CMD ["./grpc-server"]
