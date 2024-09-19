package main

import (
	"context"
	"fmt"
	"log"
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
	streaming := streaming.NewStreaming("chat", 0)
	logic := logic.NewServer(db, streaming)
	defer logic.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go logic.UploadMessages(ctx, &wg)
	wg.Wait()

}
func handleDbConnection() (*persistance.DB, error) {
	db, err := persistance.NewDB(
		"chat_db",
		"postgres",
		"chatAppPass",
		"postgres",
		"5432",
	)
	if err == nil {
		return db, nil
	}
	db, err = persistance.NewDB(
		"localhost",
		"postgres",
		"chatAppPass",
		"postgres",
		"5432",
	)
	if err == nil {
		return db, nil
	}
	return nil, fmt.Errorf("could not connect to database")
}
