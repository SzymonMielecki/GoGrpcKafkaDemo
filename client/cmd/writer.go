package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/loginState"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/userServiceClient"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/streaming/producer"
	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/types"
	pb "github.com/SzymonMielecki/GoGrpcKafkaGormDemo/usersService"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func WriterCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "writer",
		Short: "Writes messages to the chat",
		Long:  `Writes messages to the chat, you need to be logged in to use this command`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			state, err := loginState.LoadState(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if !state.LoggedIn {
				fmt.Printf("\033[1;31mYou need to be logged in to use this command\033[0m\n")
				cancel()
				return
			}
			userServiceClient, err := userServiceClient.NewUserServiceClient()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer userServiceClient.Close()
			response, err := userServiceClient.CheckUser(ctx, &pb.CheckUserRequest{
				Id:           uint32(state.Id),
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

			user := &types.User{
				Model: gorm.Model{
					ID: uint(response.User.Id),
				},
				Username:     response.User.Username,
				Email:        response.User.Email,
				PasswordHash: response.User.PasswordHash,
			}
			fmt.Println("Logged in as", user.Username)
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
			err = streaming.SendMessage(ctx, &types.Message{
				Content:  message,
				SenderID: state.Id,
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
