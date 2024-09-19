package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/SzymonMielecki/chatApp/client/cmd"
	"github.com/SzymonMielecki/chatApp/client/loginState"
	"github.com/SzymonMielecki/chatApp/client/userServiceClient"
	"github.com/SzymonMielecki/chatApp/streaming"
	"github.com/SzymonMielecki/chatApp/types"
	pb "github.com/SzymonMielecki/chatApp/usersService"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chatApp",
	Short: "ChatApp is a chat application",
	Long:  `ChatApp is a chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		readerCmd.Run(cmd, args)
	},
}

var readerCmd = &cobra.Command{
	Use:   "chatApp",
	Short: "ChatApp is a chat application",
	Long:  `ChatApp is a chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		state, err := loginState.LoadState()
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
		response, err := userServiceClient.CheckUser(context.Background(), &pb.CheckUserRequest{
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

		streaming := streaming.NewStreaming("chat", 0)
		defer streaming.Close()
		ch := make(chan *types.Message)
		var wg sync.WaitGroup
		wg.Add(1)
		go streaming.ReceiveMessages(context.Background(), ch, &wg)
		go func() {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					cancel()
					return
				case msg := <-ch:
					fmt.Printf("\033[31m%s\033[0m\n", msg.Content)
				}
			}
		}()
		wg.Wait()
		cancel()
	},
}

var writerCmd = &cobra.Command{
	Use:   "chatApp",
	Short: "ChatApp is a chat application",
	Long:  `ChatApp is a chat application`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		state, err := loginState.LoadState()
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
		response, err := userServiceClient.CheckUser(context.Background(), &pb.CheckUserRequest{
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

		streaming := streaming.NewStreaming("chat", 0)
		defer streaming.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		content := strings.Join(args, " ")
		if content[0] == '"' && content[len(content)-1] == '"' {
			content = content[1 : len(content)-1]
		}
		streaming.SendMessage(ctx, &types.Message{
			Content:  content,
			SenderID: state.Id,
		})
		cancel()
	},
}

func Execute() {
	cmd.LoginCmd.Flags().StringP("username", "u", "", "Username")
	cmd.LoginCmd.MarkFlagRequired("username")
	cmd.LoginCmd.Flags().StringP("email", "e", "", "Email")
	cmd.LoginCmd.MarkFlagRequired("email")
	cmd.RegisterCmd.Flags().StringP("password", "p", "", "Password")
	cmd.RegisterCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(cmd.RegisterCmd)

	cmd.LoginCmd.Flags().StringP("username", "u", "", "Username")
	cmd.LoginCmd.Flags().StringP("email", "e", "", "Email")
	cmd.LoginCmd.MarkFlagsOneRequired("username", "email")
	cmd.LoginCmd.Flags().StringP("password", "p", "", "Password")
	cmd.LoginCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(cmd.LoginCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
