package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/SzymonMielecki/chatApp/client/loginState"
	"github.com/SzymonMielecki/chatApp/client/userServiceClient"
	"github.com/SzymonMielecki/chatApp/streaming/client"
	"github.com/SzymonMielecki/chatApp/types"
	pb "github.com/SzymonMielecki/chatApp/usersService"
	"github.com/spf13/cobra"
)

func ReaderCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reader",
		Short: "Reader is a command that reads messages from the chat",
		Long:  `Reader is a command that reads messages from the chat`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			state, err := loginState.LoadState(ctx)
			if err != nil {
				fmt.Println("Error loading state in client/main.go: \n", err)
				cancel()
				return
			}
			userServiceClient, err := userServiceClient.NewUserServiceClient()
			if err != nil {
				fmt.Println("Error creating user service client in client/main.go: \n", err)
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
				fmt.Println("Error checking user in client/main.go: \n", err)
				cancel()
				return
			}
			if !response.Success {
				fmt.Println("Not logged in")
				cancel()
				return
			}
			fmt.Println("Logged in as", state.Username)
			streaming, err := client.NewStreamingClient(ctx, "chat", 1, []string{"localhost:9092"})
			if err != nil {
				cancel()
				return
			}
			defer streaming.Close()
			ch := make(chan *types.StreamingMessage)
			var wg sync.WaitGroup

			wg.Add(1)
			fmt.Println("Starting to receive messages")
			go streaming.ReceiveMessages(ctx, ch, &wg)

			wg.Add(1)
			fmt.Println("Starting to print messages")
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
			wg.Wait()
			cancel()
		},
	}

}
