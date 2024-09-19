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
		log.Fatalf("failed to connect to database: %v", err)
	}
	streaming := streaming.NewStreaming("kafka", "chat", 0)
	logic := logic.NewServer(db, streaming)
	defer logic.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("Starting chat server")
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
