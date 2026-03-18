package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context <agent_id> <notes>",
	Short: "Update the Middleman's notes about an agent's known context",
	Long: `Store free-text notes about what the agent knows: files read, decisions made, etc.
This is the Middleman's memory about the agent, not the agent's own context.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := orch.UpdateContext(context.Background(), args[0], args[1]); err != nil {
			return err
		}
		fmt.Printf("context updated for agent %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
