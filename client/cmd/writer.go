package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/SzymonMielecki/chatApp/client/loginState"
	"github.com/SzymonMielecki/chatApp/client/userServiceClient"
	"github.com/SzymonMielecki/chatApp/streaming/producer"
	"github.com/SzymonMielecki/chatApp/types"
	pb "github.com/SzymonMielecki/chatApp/usersService"
	"github.com/spf13/cobra"
)

func WriterCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "writer",
		Short: "Writer is a command that writes messages to the chat",
		Long:  `Writer is a command that writes messages to the chat`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			state, err := loginState.LoadState(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			userServiceClient, err := userServiceClient.NewUserServiceClient()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer userServiceClient.Close()
			response, err := userServiceClient.CheckUser(ctx, &pb.CheckUserRequest{
				Username:     state.Username,
				Email:        state.Email,
				PasswordHash: state.PasswordHash,
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if !response.Success {
				fmt.Println("Not logged in")
				os.Exit(1)
			}
			fmt.Println("Logged in as", state.Username)
			fmt.Println("Enter your message:")
			reader := bufio.NewReader(os.Stdin)
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			streaming, err := producer.NewStreamingProducer(ctx, "chat", 1, []string{"localhost:9092"})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer streaming.Close()
			var wg sync.WaitGroup
			wg.Add(1)
			err = streaming.SendMessage(ctx, &types.StreamingMessage{
				Content:        message,
				SenderID:       state.Id,
				SenderUsername: state.Username,
				SenderEmail:    state.Email,
			}, &wg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			wg.Wait()
			fmt.Println("Message sent")
			cancel()
		},
	}

}
