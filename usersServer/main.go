package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersServer/logic"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersServer/persistance"
	pb "github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersService"
	"google.golang.org/grpc"
)

func main() {
	db, err := handleDbConnection()
	if err != nil {
		log.Fatalf("failed to connect to database in usersServer/main.go: \n%v", err)
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen in usersServer/main.go: \n%v", err)
	}
	log.Printf("Server listening on port 50051")

	s := grpc.NewServer()
	logic := logic.NewServer(db)
	pb.RegisterUsersServiceServer(s, logic)
	log.Printf("Server registered and ready to serve")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve in usersServer/main.go: \n%v", err)
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
