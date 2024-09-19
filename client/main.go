package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"

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
	Use:   "reader",
	Short: "Reader is a command that reads messages from the chat",
	Long:  `Reader is a command that reads messages from the chat`,
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

		streaming := streaming.NewStreaming("kafka", "chat", 0)
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
	Use:   "writer",
	Short: "Writer is a command that writes messages to the chat",
	Long:  `Writer is a command that writes messages to the chat`,
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

		streaming := streaming.NewStreaming("kafka", "chat", 0)
		defer streaming.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		streaming.SendMessage(ctx, &types.Message{
			Content:  message,
			SenderID: state.Id,
		})
		cancel()
	},
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the chat application",
	Long:  `Login to the chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		hasher := sha256.New()
		hasher.Write([]byte(password))
		passwordHash := hex.EncodeToString(hasher.Sum(nil))
		usernameOrEmail := username
		if usernameOrEmail == "" {
			usernameOrEmail = email
		}
		user := &pb.LoginUserRequest{
			UsernameOrEmail: usernameOrEmail,
			PasswordHash:    passwordHash,
		}
		client, err := userServiceClient.NewUserServiceClient()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer client.Close()
		response, err := client.LoginUser(context.Background(), user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		state := loginState.NewLoginState(
			response.Success,
			uint(response.User.Id),
			username,
			email,
			passwordHash,
		)
		state.Save()
	},
}

var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register to the chat application",
	Long:  `Register to the chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		client, err := userServiceClient.NewUserServiceClient()
		if err != nil {
			fmt.Println(err)
			cancel()
			os.Exit(1)
			return
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
			fmt.Println(response)
			fmt.Println(err)
			cancel()
			os.Exit(1)
			return
		}
		state := loginState.NewLoginState(
			response.Success,
			uint(response.User.Id),
			username,
			email,
			passwordHash,
		)
		state.Save()
		cancel()
	},
}

var username string
var email string
var password string
var message string

func Execute() {
	// Add flags for LoginCmd
	LoginCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	LoginCmd.Flags().StringVarP(&email, "email", "e", "", "Email")
	LoginCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	LoginCmd.MarkFlagsOneRequired("username", "email")
	LoginCmd.MarkFlagRequired("password")

	// Add flags for RegisterCmd
	RegisterCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	RegisterCmd.Flags().StringVarP(&email, "email", "e", "", "Email")
	RegisterCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	RegisterCmd.MarkFlagRequired("username")
	RegisterCmd.MarkFlagRequired("email")
	RegisterCmd.MarkFlagRequired("password")

	// Add flags for writerCmd
	writerCmd.Flags().StringVarP(&message, "message", "m", "", "Message")
	writerCmd.MarkFlagRequired("message")

	// Add commands to rootCmd
	rootCmd.AddCommand(LoginCmd)
	rootCmd.AddCommand(RegisterCmd)
	rootCmd.AddCommand(writerCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
