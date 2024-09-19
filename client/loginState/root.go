package loginState

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	pb "github.com/SzymonMielecki/chatApp/usersService"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginState struct {
	LoggedIn     bool   `json:"logged_in"`
	Id           uint   `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
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

func LoadState() (*LoginState, error) {
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

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	c := pb.NewUsersServiceClient(conn)
	response, err := c.CheckUser(context.Background(), &pb.CheckUserRequest{
		Username:     s.Username,
		Email:        s.Email,
		PasswordHash: s.PasswordHash,
	})
	if err != nil {
		return nil, err
	}
	if response.Success {
		s.LoggedIn = true
		s.Id = uint(response.User.Id)
	}
	return nil, fmt.Errorf("invalid credentials")
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
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}
}
