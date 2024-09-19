package main

import (
	"fmt"
	"log"
	"net"
	"os"

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
	pb.RegisterUsersServiceServer(s, logic)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

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
