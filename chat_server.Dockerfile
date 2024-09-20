FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

# Keep CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o chat_server chatServer/main.go

# Use Ubuntu 22.04 as the base image
FROM ubuntu:22.04

# Install necessary libraries
RUN apt-get update && apt-get install -y \
    ca-certificates \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/chat_server /chat_server

CMD ["./chat_server"]
