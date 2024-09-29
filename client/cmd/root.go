package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/SzymonMielecki/GoGrpcKafkaGormDemo/client/loginState"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "chatUp",
		Short: "ChatUp is a chat application",
		Long:  `ChatUp is a real-time chat application based on Kafka and gRPC`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			state, err := loginState.LoadState(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Welcome to ChatUp!")
			if state.LoggedIn {
				fmt.Println("You are logged in")
			} else {
				fmt.Println("You are not logged in")
			}
			fmt.Println("For a list of commands, type 'help'")
			cancel()
		},
	}
}
