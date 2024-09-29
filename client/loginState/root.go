package loginState

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	pb "github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersService"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginState struct {
	LoggedIn     bool   `json:"logged_in"`
	Id           uint   `json:"id"`
	PasswordHash string `json:"password_hash"`
}

func (s *LoginState) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	stateFile := filepath.Join(homeDir, ".chatapp_session")
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(stateFile, []byte(data), 0644)
}

func LoadState(ctx context.Context) (*LoginState, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	stateFile := filepath.Join(homeDir, ".chatapp_session")
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return nil, err
	}

	s := &LoginState{}
	err = json.Unmarshal(data, s)
	if err != nil {
		return nil, err
	}
	if !s.LoggedIn {
		return s, fmt.Errorf("you need to be logged in to use this command")
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer conn.Close()
	c := pb.NewUsersServiceClient(conn)
	response, err := c.CheckUser(ctx, &pb.CheckUserRequest{
		PasswordHash: s.PasswordHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check user in loginState/root.go: \n%v", err)
	}
	if response.Success {
		s.LoggedIn = true
		s.Id = uint(response.User.Id)
		return s, nil
	}
	fmt.Println("User not found")
	return nil, fmt.Errorf("user not found")
}

func NewLoginState(
	success bool,
	id uint,
	username string,
	email string,
	passwordHash string,
) *LoginState {
	return &LoginState{
		LoggedIn:     success,
		Id:           id,
		PasswordHash: passwordHash,
	}
}
