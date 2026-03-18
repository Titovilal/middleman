package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var rewindFlags struct {
	to   string
	list bool
}

var rewindCmd = &cobra.Command{
	Use:   "rewind <agent_id>",
	Short: "Rewind an agent to a checkpoint",
	Long: `Fork an agent back to a previous checkpoint. The original session is preserved.
Default (no --to) rewinds to the most recent checkpoint (undoes last task).
Use --list to show available checkpoints without rewinding.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentID := args[0]

		if rewindFlags.list {
			a, err := orch.Inspect(context.Background(), agentID)
			if err != nil {
				return err
			}
			if len(a.Checkpoints) == 0 {
				fmt.Printf("no checkpoints for agent %s\n", agentID)
				return nil
			}
			fmt.Printf("checkpoints for %s:\n", agentID)
			for i, cp := range a.Checkpoints {
				fmt.Printf("  [%d] %-40s turn=%-3d  %s\n", i, cp.Label, cp.TurnIndex, cp.CreatedAt.Format("2006-01-02 15:04:05"))
			}
			return nil
		}

		a, err := orch.Rewind(context.Background(), agentID, rewindFlags.to)
		if err != nil {
			return err
		}

		cp := a.LatestCheckpoint()
		label := "unknown"
		if cp != nil {
			label = cp.Label
		}
		fmt.Printf("rewound %s to checkpoint %q\n", agentID, label)
		fmt.Printf("  new session: %s\n", a.SessionID)
		return nil
	},
}

func init() {
	rewindCmd.Flags().StringVar(&rewindFlags.to, "to", "", "checkpoint label to rewind to (default: latest)")
	rewindCmd.Flags().BoolVar(&rewindFlags.list, "list", false, "list checkpoints without rewinding")
	rootCmd.AddCommand(rewindCmd)
}
