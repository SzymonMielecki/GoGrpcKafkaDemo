FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o users_server usersServer/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates postgresql-client 

COPY --from=builder /app/users_server /users_server

EXPOSE 50051

CMD ["./users_server"]
