# ChatUp - GoGrpcKafkaDemo

This is a chat app based on microservices architecture.

## Technologies used

-   [x] Go
-   [x] Docker
-   [x] gRPC
-   [x] Kafka
-   [x] PostgresQl
-   [x] gORM
-   [x] Redis

## Usage Details

### Run the application

```bash
docker compose up --build
```

### Run the client

```bash
go build -o chatUp client/main.go
./chatUp
```

```
ChatUp is a real-time chat application based on Kafka and gRPC

Usage:
  chatUp [flags]
  chatUp [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  login       Login to the chat application
  reader      Reads messages from the chat
  register    Register to the chat application
  writer      Writes messages to the chat

Flags:
  -h, --help   help for chatApp

Use "chatUp [command] --help" for more information about a command.
```
