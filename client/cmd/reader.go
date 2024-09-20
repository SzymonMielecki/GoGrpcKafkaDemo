package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/loginState"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/userServiceClient"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/streaming/client"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/types"
	pb "github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersService"
	"github.com/spf13/cobra"
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
				Username:     state.Username,
				Email:        state.Email,
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
			fmt.Println("Logged in as", state.Username)
			streaming, err := client.NewStreamingClient(ctx, "chat", 1, []string{"localhost:9092"})
			if err != nil {
				fmt.Printf("\033[1;31mFailed to create streaming client in client/cmd/reader.go: \n%v\033[0m", err)
				cancel()
				return
			}
			defer streaming.Close()
			ch := make(chan *types.StreamingMessage)
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
						fmt.Printf("\033[1;32m%d:\033[0m %s", msg.SenderID, msg.Content)
					}
				}
			}()
			fmt.Println("The chat is running, press Ctrl+C to stop")
			wg.Wait()
			cancel()
		},
	}

}
