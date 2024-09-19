package userServiceClient

import (
	"context"
	"fmt"

	pb "github.com/SzymonMielecki/chatApp/usersService"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb.UsersServiceClient
	conn *grpc.ClientConn
}

func NewUserServiceClient() (*Client, error) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to users_server in userServiceClient/root.go: \n%v", err)
	}
	c := pb.NewUsersServiceClient(conn)
	return &Client{conn: conn, UsersServiceClient: c}, nil
}
func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) LoginUser(ctx context.Context, user *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	r, err := c.UsersServiceClient.LoginUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to login user in userServiceClient/root.go: \n%v", err)
	}
	return r, nil
}

func (c *Client) RegisterUser(ctx context.Context, user *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	r, err := c.UsersServiceClient.RegisterUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to register user in userServiceClient/root.go: \n%v", err)
	}
	return r, nil
}
func (c *Client) GetUser(ctx context.Context, user *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return c.UsersServiceClient.GetUser(ctx, user)
}
func (c *Client) CheckUser(ctx context.Context, user *pb.CheckUserRequest) (*pb.CheckUserResponse, error) {
	r, err := c.UsersServiceClient.CheckUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to check user in userServiceClient/root.go: \n%v", err)
	}
	return r, nil
}
