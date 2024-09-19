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
	"github.com/SzymonMielecki/chatApp/streaming/consumer"
	"github.com/SzymonMielecki/chatApp/streaming/producer"
	"github.com/SzymonMielecki/chatApp/types"
	pb "github.com/SzymonMielecki/chatApp/usersService"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chatApp",
	Short: "ChatApp is a chat application",
	Long:  `ChatApp is a chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		state, err := loginState.LoadState(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Welcome to ChatApp!")
		fmt.Println("You are logged in as", state.Username)
		fmt.Println("You can use the following commands:")
		fmt.Println("reader - reads messages from the chat")
		fmt.Println("writer - writes messages to the chat")
		fmt.Println("login - logs you in")
		fmt.Println("register - registers you to the chat")
		cancel()
	},
}

var readerCmd = &cobra.Command{
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
		streaming, err := consumer.NewStreamingConsumer(ctx, "localhost:9092", 0, []string{"localhost:9092"})
		if err != nil {
			cancel()
			return
		}
		defer streaming.Close()
		ch := make(chan *types.Message)
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

		streaming, err := producer.NewStreamingProducer(ctx, "localhost:9092", 0, []string{"localhost:9092"})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer streaming.Close()
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

var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register to the chat application",
	Long:  `Register to the chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		client, err := userServiceClient.NewUserServiceClient()
		if err != nil {
			fmt.Printf("create user in client/main.go: \n%v", err)
			os.Exit(1)
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
			fmt.Printf("Failed to register user in client/main.go: \n%v", err)
			os.Exit(1)
		}
		state := loginState.NewLoginState(
			response.Success,
			uint(response.User.Id),
			username,
			email,
			passwordHash,
		)
		defer state.Save()
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
	rootCmd.AddCommand(readerCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
