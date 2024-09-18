package cmd

import (
	"context"
	"fmt"
	"os"

	pb "github.com/SzymonMielecki/chatApp/usersService"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var rootCmd = &cobra.Command{
	Use:   "chatApp",
	Short: "ChatApp is a chat application",
	Long:  `ChatApp is a chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root command")
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the chat application",
	Long:  `Login to the chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Login command")
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register to the chat application",
	Long:  `Register to the chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer conn.Close()
		c := pb.NewUsersServiceClient(conn)
		username, _ := cmd.Flags().GetString("username")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		user := &pb.RegisterUserRequest{
			Username: username,
			Email:    email,
			Password: password,
		}
		response, err := c.RegisterUser(context.Background(), user)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(response)
	},
}

func Execute() {
	registerCmd.Flags().StringP("username", "u", "", "Username")
	registerCmd.MarkFlagRequired("username")
	registerCmd.Flags().StringP("email", "e", "", "Email")
	registerCmd.MarkFlagRequired("email")
	registerCmd.Flags().StringP("password", "p", "", "Password")
	registerCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(registerCmd)

	loginCmd.Flags().StringP("username", "u", "", "Username")
	loginCmd.Flags().StringP("email", "e", "", "Email")
	loginCmd.MarkFlagsOneRequired("username", "email")
	loginCmd.Flags().StringP("password", "p", "", "Password")
	loginCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(loginCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
