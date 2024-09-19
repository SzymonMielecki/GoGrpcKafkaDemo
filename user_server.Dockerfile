FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o main usersServer/main.go

COPY wait-for-db.sh /wait-for-db.sh

FROM alpine:latest

RUN apk add --no-cache ca-certificates postgresql-client 
RUN apk add --no-cache netcat-openbsd 

COPY --from=builder /app/wait-for-db.sh /wait-for-db.sh
COPY --from=builder /app/main /main


RUN chmod +x /wait-for-db.sh

EXPOSE 50051

CMD ["/wait-for-db.sh", "users_db", "5432", "./main"]
