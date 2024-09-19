package login

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	pb "github.com/SzymonMielecki/chatApp/usersService"

	"github.com/SzymonMielecki/chatApp/client/state"
	"github.com/SzymonMielecki/chatApp/client/userServiceClient"
	"github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
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
		client, err := userServiceClient.NewUserServiceClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer client.Close()
		response, err := client.LoginUser(context.Background(), user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		state := &state.LoginState{
			LoggedIn:     response.Success,
			Id:           uint(response.User.Id),
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		}
		state.Save()
	},
}

var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register to the chat application",
	Long:  `Register to the chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := userServiceClient.NewUserServiceClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer client.Close()
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
		response, err := client.RegisterUser(context.Background(), user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		state := &state.LoginState{
			LoggedIn:     response.Success,
			Id:           uint(response.User.Id),
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		}
		state.Save()
	},
}
