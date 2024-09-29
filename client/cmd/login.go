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

func LoginCommand(username, email, password *string) *cobra.Command {
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
			hasher.Write([]byte(*password))
			passwordHash := hex.EncodeToString(hasher.Sum(nil))
			user := &pb.LoginUserRequest{
				Username:     *username,
				Email:        *email,
				PasswordHash: passwordHash,
			}
			response, err := client.LoginUser(context.Background(), user)
			if err != nil {
				fmt.Println(err)
				return
			}
			login_state := loginState.NewLoginState(
				response.Success,
				uint(response.User.Id),
				response.User.Username,
				response.User.Email,
				response.User.PasswordHash,
			)
			login_state.Save()
			tagline := strings.Split(*email, "@")[0]
			color := utils.GetColorForUser(*username)
			if response.Success {
				fmt.Printf("Successfully logged in as \033[%dm%s@%s\033[0m\n", color, *username, tagline)
			} else {
				fmt.Printf("\033[1;31mRegistration failed: %s\033[0m\n", response.Message)
			}
		},
	}

}
