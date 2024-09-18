FROM golang:1.22 as builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config.yaml .

EXPOSE 50051

CMD ["./main"]