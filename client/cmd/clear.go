package cmd

import (
	"context"
	"fmt"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/client/loginState"
	"github.com/spf13/cobra"
)

func ClearCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Clears the current session",
		Long:  `Clears the current session`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			state, err := loginState.LoadState(ctx)
			if err != nil {
				return
			}
			state.Clear()
			state.Save()
			fmt.Println("Session cleared")
		},
	}
}
