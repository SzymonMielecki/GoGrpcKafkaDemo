package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	pb "github.com/SzymonMielecki/chatApp/usersService"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"crypto/sha256"
)

type LoginState struct {
	LoggedIn     bool   `json:"logged_in"`
	Id           string `json:"id"`
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
		s.Id = response.Id
	}
	return nil, fmt.Errorf("invalid credentials")
}

var rootCmd = &cobra.Command{
	Use:   "chatApp",
	Short: "ChatApp is a chat application",
	Long:  `ChatApp is a chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root command")
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the chat application",
	Long:  `Login to the chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		hasher := sha256.New()
		hasher.Write([]byte(password))
		passwordHash := hex.EncodeToString(hasher.Sum(nil))
		usernameOrEmail := username
		if usernameOrEmail == "" {
			usernameOrEmail = email
		}
		user := &pb.LoginUserRequest{
			UsernameOrEmail: usernameOrEmail,
			PasswordHash:    passwordHash,
		}
		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer conn.Close()
		c := pb.NewUsersServiceClient(conn)
		response, err := c.LoginUser(context.Background(), user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		state := &LoginState{
			LoggedIn:     response.Success,
			Id:           response.Id,
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		}
		state.Save()
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register to the chat application",
	Long:  `Register to the chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer conn.Close()
		c := pb.NewUsersServiceClient(conn)
		username, _ := cmd.Flags().GetString("username")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		hasher := sha256.New()
		hasher.Write([]byte(password))
		passwordHash := hex.EncodeToString(hasher.Sum(nil))
		user := &pb.RegisterUserRequest{
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		}
		response, err := c.RegisterUser(context.Background(), user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		state := &LoginState{
			LoggedIn:     response.Success,
			Id:           response.Id,
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		}
		state.Save()
	},
}

func Execute() {
	registerCmd.Flags().StringP("username", "u", "", "Username")
	registerCmd.MarkFlagRequired("username")
	registerCmd.Flags().StringP("email", "e", "", "Email")
	registerCmd.MarkFlagRequired("email")
	registerCmd.Flags().StringP("password", "p", "", "Password")
	registerCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(registerCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("email", "e", "", "Email")
	loginCmd.MarkFlagsOneRequired("username", "email")
	loginCmd.Flags().StringP("password", "p", "", "Password")
	loginCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(loginCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
