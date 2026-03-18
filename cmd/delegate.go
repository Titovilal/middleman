package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Titovilal/middleman/agent"
	"github.com/spf13/cobra"
)

var delegateFlags struct {
	to      string
	timeout time.Duration
}

var delegateCmd = &cobra.Command{
	Use:   "delegate <prompt>",
	Short: "Send a task to an agent (runs in background)",
	Long:  `Send a task to a specific agent (--to). The task runs in a background process. If the agent is busy, the task is queued and runs automatically when the current task finishes.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if delegateFlags.to == "" {
			return fmt.Errorf("--to <agent_id> is required")
		}

		prompt := strings.Join(args, " ")

		task, err := orch.DelegateAsync(context.Background(), delegateFlags.to, prompt, delegateFlags.timeout)
		if err != nil {
			return err
		}

		// Only launch a background process if the task is pending (not queued).
		// Queued tasks are picked up automatically by the already-running background process.
		if task.Status == agent.TaskPending {
			exe, err := os.Executable()
			if err != nil {
				return fmt.Errorf("resolve executable: %w", err)
			}

			bgArgs := []string{"_run-task", delegateFlags.to, task.TaskID}
			if delegateFlags.timeout > 0 {
				bgArgs = append(bgArgs, "--timeout", delegateFlags.timeout.String())
			}

			bgCmd := exec.Command(exe, bgArgs...)
			bgCmd.Stdout = nil
			bgCmd.Stderr = nil
			bgCmd.Stdin = nil
			if err := bgCmd.Start(); err != nil {
				return fmt.Errorf("start background task: %w", err)
			}
			_ = bgCmd.Process.Release()
		}

		fmt.Printf("task delegated\n")
		fmt.Printf("  agent:   %s\n", delegateFlags.to)
		fmt.Printf("  task_id: %s\n", task.TaskID)
		fmt.Printf("  status:  %s\n", task.Status)
		if task.Status == agent.TaskQueued {
			fmt.Printf("\nAgent is busy. Task queued — it will run when the current task finishes.\n")
		}
		fmt.Printf("\nCheck result:\n")
		fmt.Printf("  mdm result %s\n", delegateFlags.to)
		return nil
	},
}

func init() {
	delegateCmd.Flags().StringVar(&delegateFlags.to, "to", "", "agent ID to delegate to (required)")
	delegateCmd.Flags().DurationVar(&delegateFlags.timeout, "timeout", 0, "task timeout (default 5m, e.g. 10m, 2m30s)")
	rootCmd.AddCommand(delegateCmd)
}
