package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/SzymonMielecki/chatApp/chatServer/logic"
	"github.com/SzymonMielecki/chatApp/chatServer/persistance"
	"github.com/SzymonMielecki/chatApp/streaming"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := handleDbConnection()
	if err != nil {
		log.Fatalf("failed to connect to database in chatServer/main.go: \n%v", err)
	}
	broker := os.Getenv("KAFKA_BROKER")
	streaming, err := streaming.NewStreaming(ctx, broker, "chat", 0, []string{broker})
	if err != nil {
		log.Fatalf("failed to create streaming in chatServer/main.go: \n%v", err)
	}
	defer streaming.Close()
	logic := logic.NewServer(db, streaming)
	defer logic.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	log.Default().Println("Starting chat server")
	go logic.UploadMessages(ctx, &wg)
	wg.Wait()
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
