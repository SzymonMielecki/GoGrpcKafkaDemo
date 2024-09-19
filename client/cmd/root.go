package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/SzymonMielecki/chatApp/client/cmd/login"
	"github.com/SzymonMielecki/chatApp/client/state"
	"github.com/SzymonMielecki/chatApp/client/userServiceClient"
	pb "github.com/SzymonMielecki/chatApp/usersService"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chatApp",
	Short: "ChatApp is a chat application",
	Long:  `ChatApp is a chat application`,
	Run: func(cmd *cobra.Command, args []string) {
		state, err := state.LoadState()
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
	},
}

func Execute() {
	login.LoginCmd.Flags().StringP("username", "u", "", "Username")
	login.LoginCmd.MarkFlagRequired("username")
	login.LoginCmd.Flags().StringP("email", "e", "", "Email")
	login.RegisterCmd.MarkFlagRequired("email")
	login.RegisterCmd.Flags().StringP("password", "p", "", "Password")
	login.RegisterCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(login.RegisterCmd)

	login.LoginCmd.Flags().StringP("username", "u", "", "Username")
	login.LoginCmd.Flags().StringP("email", "e", "", "Email")
	login.LoginCmd.MarkFlagsOneRequired("username", "email")
	login.LoginCmd.Flags().StringP("password", "p", "", "Password")
	login.LoginCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(login.LoginCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
