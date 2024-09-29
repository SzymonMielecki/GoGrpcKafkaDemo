package cmd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/loginState"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/userServiceClient"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/utils"
	pb "github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersService"
	"github.com/spf13/cobra"
)

func RegisterCommand(username, email, password *string) *cobra.Command {
	return &cobra.Command{
		Use:   "register",
		Short: "Register to the chat application",
		Long:  `Register to the chat application`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client, err := userServiceClient.NewUserServiceClient()
			if err != nil {
				fmt.Printf("\033[1;31mFailed to create user in client/main.go: \n%v\033[0m", err)
				os.Exit(1)
			}
			defer client.Close()

			hasher := sha256.New()
			hasher.Write([]byte(*password))
			passwordHash := hex.EncodeToString(hasher.Sum(nil))
			user := &pb.RegisterUserRequest{
				Username:     *username,
				Email:        *email,
				PasswordHash: passwordHash,
			}
			response, err := client.RegisterUser(ctx, user)
			if err != nil {
				fmt.Printf("\033[1;31mFailed to register user in client/main.go: \n%v\033[0m", err)
				os.Exit(1)
			}
			state := loginState.NewLoginState(
				response.Success,
				uint(response.User.Id),
				*username,
				*email,
				passwordHash,
			)
			state.Save()
			tagline := strings.Split(*email, "@")[0]
			color := utils.GetColorForUser(*username)
			if response.Success {
				fmt.Printf("Successfully registered as \033[%dm%s@%s\033[0m\n", color, *username, tagline)
			} else {
				fmt.Printf("\033[1;31mRegistration failed: %s\033[0m\n", response.Message)
			}
		},
	}
}
