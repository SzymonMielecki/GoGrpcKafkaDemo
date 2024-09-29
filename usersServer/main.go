package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersServer/logic"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersServer/persistance"
	pb "github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersService"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	db, c, err := handleConnections()
	if err != nil {
		log.Fatalf("failed to connect to database in usersServer/main.go: \n%v", err)
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen in usersServer/main.go: \n%v", err)
	}
	log.Printf("Server listening on port 50051")

	s := grpc.NewServer()
	logic := logic.NewServer(db, c)
	pb.RegisterUsersServiceServer(s, logic)
	log.Printf("Server registered and ready to serve")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve in usersServer/main.go: \n%v", err)
	}

}

func handleConnections() (*persistance.DB, *cache.Cache, error) {
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
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database connection in usersServer/main.go: \n%v", err)
	}
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
	})
	c := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, 0),
	})
	return db, c, nil
}
