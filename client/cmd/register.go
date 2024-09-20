package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/SzymonMielecki/chatApp/client/loginState"
	"github.com/SzymonMielecki/chatApp/client/userServiceClient"
	pb "github.com/SzymonMielecki/chatApp/usersService"
	"github.com/spf13/cobra"
)

func RegisterCommand(username, email, password string) *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Register to the chat application",
		Long:  `Register to the chat application`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client, err := userServiceClient.NewUserServiceClient()
			if err != nil {
				fmt.Printf("create user in client/main.go: \n%v", err)
				os.Exit(1)
			}
			defer client.Close()

			hasher := sha256.New()
			hasher.Write([]byte(password))
			passwordHash := hex.EncodeToString(hasher.Sum(nil))
			user := &pb.RegisterUserRequest{
				Username:     username,
				Email:        email,
				PasswordHash: passwordHash,
			}
			response, err := client.RegisterUser(ctx, user)
			if err != nil {
				fmt.Printf("Failed to register user in client/main.go: \n%v", err)
				os.Exit(1)
			}
			state := loginState.NewLoginState(
				response.Success,
				uint(response.User.Id),
				username,
				email,
				passwordHash,
			)
			defer state.Save()
		},
	}
}
