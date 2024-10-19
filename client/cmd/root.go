package cmd

import (
	"context"
	"fmt"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/client/loginState"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "chatUp",
		Short: "ChatUp is a chat application",
		Long:  `ChatUp is a real-time chat application based on Kafka and gRPC`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			state, err := loginState.LoadState(ctx)
			fmt.Println("Welcome to ChatUp!")
			if err != nil || !state.LoggedIn {
				fmt.Println("You are not logged in")
			} else {
				fmt.Println("You are logged in")
			}
			fmt.Println("For a list of commands, type 'help'")
		},
	}
}
