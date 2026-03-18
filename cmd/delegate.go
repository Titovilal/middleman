package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var delegateFlags struct {
	to      string
	timeout time.Duration
}

var delegateCmd = &cobra.Command{
	Use:   "delegate <prompt>",
	Short: "Send a task to an agent",
	Long:  `Send a task to a specific agent (--to). The full response is printed to stdout. Internal agent work is not shown.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if delegateFlags.to == "" {
			return fmt.Errorf("--to <agent_id> is required")
		}

		prompt := strings.Join(args, " ")

		task, err := orch.Delegate(context.Background(), delegateFlags.to, prompt, delegateFlags.timeout)
		if err != nil {
			return err
		}

		if task.IsError {
			fmt.Printf("[error] %s\n", task.ErrorDetail)
			return nil
		}

		fmt.Println(task.Response)
		return nil
	},
}

func init() {
	delegateCmd.Flags().StringVar(&delegateFlags.to, "to", "", "agent ID to delegate to (required)")
	delegateCmd.Flags().DurationVar(&delegateFlags.timeout, "timeout", 0, "task timeout (default 5m, e.g. 10m, 2m30s)")
	rootCmd.AddCommand(delegateCmd)
}
