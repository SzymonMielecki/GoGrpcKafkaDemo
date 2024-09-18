package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SzymonMielecki/chatApp/usersServer/logic"
	"github.com/SzymonMielecki/chatApp/usersServer/persistance"
	pb "github.com/SzymonMielecki/chatApp/usersService"
	"google.golang.org/grpc"
)

func main() {
	db, err := handleDbConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	logic := logic.NewServer(db)
	pb.RegisterUsersServiceServer(s, logic.UnimplementedUsersServiceServer)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func handleDbConnection() (*persistance.DB, error) {
	db, err := persistance.NewDB(
		"db",
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
