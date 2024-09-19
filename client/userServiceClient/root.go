package userServiceClient

import (
	"context"
	"fmt"
	"os"

	pb "github.com/SzymonMielecki/chatApp/usersService"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	conn   *grpc.ClientConn
	client pb.UsersServiceClient
}

func NewUserServiceClient() (*UserServiceClient, error) {
	conn, err := newUsersConn()
	if err != nil {
		return nil, err
	}
	return &UserServiceClient{conn: conn, client: pb.NewUsersServiceClient(conn)}, nil
}
func (c *UserServiceClient) Close() {
	c.conn.Close()
}

func newUsersConn() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient("usersServer:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		conn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return conn, nil
}

func (c *UserServiceClient) LoginUser(ctx context.Context, user *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	return c.client.LoginUser(ctx, user)
}

func (c *UserServiceClient) RegisterUser(ctx context.Context, user *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	return c.client.RegisterUser(ctx, user)
}
func (c *UserServiceClient) GetUser(ctx context.Context, user *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return c.client.GetUser(ctx, user)
}
func (c *UserServiceClient) CheckUser(ctx context.Context, user *pb.CheckUserRequest) (*pb.CheckUserResponse, error) {
	return c.client.CheckUser(ctx, user)
}
