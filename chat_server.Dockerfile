FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o main chatServer/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates postgresql-client 
RUN apk add --no-cache netcat-openbsd 
RUN apk add --no-cache librdkafka-dev pkgconf
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY --from=builder /app/wait-for-kafka-and-db.sh .
COPY --from=builder /app/main .

RUN chmod +x /app/wait-for-kafka-and-db.sh


CMD ["/app/wait-for-kafka-and-db.sh", "chat_db", "5432", "kafka", "9092", "./main"]
