package cmd

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/client/loginState"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/client/userServiceClient"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/client/utils"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/streaming/client"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	pb "github.com/SzymonMielecki/GoGrpcKafkaDemo/usersService"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func ReaderCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reader",
		Short: "Reads messages from the chat",
		Long:  `Reads messages from the chat, you need to be logged in to use this command`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			state, err := loginState.LoadState(ctx)
			if err != nil {
				cancel()
				return
			}
			if !state.LoggedIn {
				fmt.Printf("\033[1;31mYou need to be logged in to use this command\033[0m\n")
				cancel()
				return
			}
			userServiceClient, err := userServiceClient.NewUserServiceClient()
			if err != nil {
				fmt.Printf("\033[1;31mFailed to create user service client in client/cmd/reader.go: \n%v\033[0m", err)
				cancel()
				return
			}
			defer userServiceClient.Close()
			response, err := userServiceClient.CheckUser(ctx, &pb.CheckUserRequest{
				Id:           uint32(state.Id),
				PasswordHash: state.PasswordHash,
			})
			if err != nil {
				fmt.Printf("\033[1;31mFailed to check user in client/cmd/reader.go: \n%v\033[0m", err)
				cancel()
				return
			}
			if !response.Success {
				fmt.Printf("\033[1;31mNot logged in\033[0m")
				cancel()
				return
			}
			user := &types.User{
				Model: gorm.Model{
					ID: uint(response.User.Id),
				},
				Username: response.User.Username,
				Email:    response.User.Email,
			}
			tagline := strings.Split(user.Email, "@")[0]
			color := utils.GetColorForUser(user.Username)
			fmt.Printf("Logged in as \033[%dm%s@%s\033[0m\n", color, user.Username, tagline)
			streaming, err := client.NewStreamingClient(ctx, "chat", 1, []string{"localhost:9092"})
			if err != nil {
				fmt.Printf("\033[1;31mFailed to create streaming client in client/cmd/reader.go: \n%v\033[0m", err)
				cancel()
				return
			}
			defer streaming.Close()
			ch := make(chan *types.Message)
			var wg sync.WaitGroup

			wg.Add(1)

			go streaming.ReceiveMessages(ctx, ch, &wg)

			wg.Add(1)

			go func() {
				for {
					select {
					case <-ctx.Done():
						wg.Done()
						cancel()
						return
					case msg := <-ch:
						sender, err := userServiceClient.GetUser(ctx, &pb.GetUserRequest{
							Id: uint32(msg.SenderID)})
						if err != nil {
							fmt.Printf("\033[1;31mFailed to get user in client/cmd/reader.go: \n%v\033[0m", err)
						}
						tagline := strings.Split(sender.User.Email, "@")[0]
						color := utils.GetColorForUser(sender.User.Username)
						fmt.Printf("\033[%dm%s@%s:\033[0m %s\n", color, sender.User.Username, tagline, msg.Content)
					}
				}
			}()
			fmt.Println("The chat is running, press Ctrl+C to stop")
			wg.Wait()
			cancel()
		},
	}

}
