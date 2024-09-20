package main

import (
	"fmt"
	"os"

	"github.com/SzymonMielecki/chatApp/client/cmd"
)

var username string
var email string
var password string

func Execute() {
	RootCmd := cmd.RootCommand()
	ReaderCmd := cmd.ReaderCommand()
	RegisterCmd := cmd.RegisterCommand(username, email, password)
	LoginCmd := cmd.LoginCommand(username, email, password)
	WriterCmd := cmd.WriterCommand()
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

	// Add commands to rootCmd
	RootCmd.AddCommand(LoginCmd)
	RootCmd.AddCommand(RegisterCmd)
	RootCmd.AddCommand(WriterCmd)
	RootCmd.AddCommand(ReaderCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
