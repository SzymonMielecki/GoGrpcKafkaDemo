package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SzymonMielecki/chatApp/chatServer/logic"
	"github.com/SzymonMielecki/chatApp/chatServer/persistance"
	"github.com/SzymonMielecki/chatApp/chatServer/streaming"
	pb "github.com/SzymonMielecki/chatApp/chatService"
	"google.golang.org/grpc"
)

func main() {
	db, err := handleDbConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	streaming := streaming.NewStreaming("chat", 0)
	s := grpc.NewServer()
	logic := logic.NewServer(db, streaming)
	defer logic.Close()
	pb.RegisterChatServiceServer(s, logic)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
