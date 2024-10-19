package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/chatServer/logic"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/chatServer/persistance"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/streaming/client"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := handleDbConnection()
	if err != nil {
		log.Fatalf("failed to connect to database in chatServer/main.go: \n%v", err)
	}
	streaming, err := client.NewStreamingClient(ctx, "chat", 0, []string{"kafka:9092", "kafka:29092"})
	if err != nil {
		log.Fatalf("failed to create streaming in chatServer/main.go: \n%v", err)
	}
	defer streaming.Close()
	logic := logic.NewServer(db, streaming)
	defer logic.Close()
	log.Default().Println("Starting chat server")
	logic.UploadMessages(ctx)
}

func handleDbConnection() (*persistance.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err := persistance.NewDB(
		dbHost,
		dbName,
		dbPassword,
		dbUser,
		dbPort,
	)
	if err == nil {
		return db, nil
	}
	return nil, fmt.Errorf("could not connect to database")
}
