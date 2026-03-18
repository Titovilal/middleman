package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var discardCmd = &cobra.Command{
	Use:   "discard <agent_id>",
	Short: "Mark an agent as discarded",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := orch.Discard(context.Background(), args[0]); err != nil {
			return err
		}
		fmt.Printf("agent %s discarded\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(discardCmd)
}
