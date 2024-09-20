package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/loginState"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/userServiceClient"
	pb "github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersService"
	"github.com/spf13/cobra"
)

func LoginCommand(username, email, password string) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to the chat application",
		Long:  `Login to the chat application`,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := userServiceClient.NewUserServiceClient()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer client.Close()
			hasher := sha256.New()
			hasher.Write([]byte(password))
			passwordHash := hex.EncodeToString(hasher.Sum(nil))
			user := &pb.LoginUserRequest{
				Username:     username,
				Email:        email,
				PasswordHash: passwordHash,
			}
			login_response, err := client.LoginUser(context.Background(), user)
			if err != nil {
				fmt.Println(err)
				return
			}
			login_state := loginState.NewLoginState(
				login_response.Success,
				uint(login_response.User.Id),
				login_response.User.Username,
				login_response.User.Email,
				login_response.User.PasswordHash,
			)
			login_state.Save()
			fmt.Println(login_response.Message)
		},
	}

}
