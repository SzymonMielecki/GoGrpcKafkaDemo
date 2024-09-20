package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/SzymonMielecki/chatApp/client/loginState"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	return &cobra.Command{
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
}
