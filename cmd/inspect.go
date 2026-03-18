package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var inspectFlags struct {
	jsonOut bool
}

var inspectCmd = &cobra.Command{
	Use:   "inspect <agent_id>",
	Short: "Show full agent details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := orch.Inspect(context.Background(), args[0])
		if err != nil {
			return err
		}

		if inspectFlags.jsonOut {
			data, _ := json.MarshalIndent(a, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("id:           %s\n", a.ID)
		fmt.Printf("connector:    %s\n", a.ConnectorName)
		fmt.Printf("session:      %s\n", a.SessionID)
		fmt.Printf("workdir:      %s\n", a.WorkDir)
		fmt.Printf("status:       %s\n", a.Status)
		fmt.Printf("created:      %s\n", a.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("last active:  %s\n", a.LastActiveAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("tasks:        %d\n", len(a.TaskLog))
		fmt.Printf("checkpoints:  %d\n", len(a.Checkpoints))

		if a.Briefing != "" {
			fmt.Printf("\nbriefing:\n  %s\n", a.Briefing)
		}
		if a.KnownContext != "" {
			fmt.Printf("\nknown context:\n  %s\n", a.KnownContext)
		}

		if len(a.Checkpoints) > 0 {
			fmt.Printf("\ncheckpoints:\n")
			for i, cp := range a.Checkpoints {
				fmt.Printf("  [%d] %-40s turn=%-3d  %s\n",
					i, cp.Label, cp.TurnIndex, cp.CreatedAt.Format("2006-01-02 15:04:05"))
			}
		}
		return nil
	},
}

func init() {
	inspectCmd.Flags().BoolVar(&inspectFlags.jsonOut, "json", false, "output as JSON")
	rootCmd.AddCommand(inspectCmd)
}
