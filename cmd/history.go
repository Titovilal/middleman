package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var historyFlags struct {
	tail int
}

var historyCmd = &cobra.Command{
	Use:   "history <agent_id>",
	Short: "Show task history for an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := orch.Inspect(context.Background(), args[0])
		if err != nil {
			return err
		}

		log := a.TaskLog
		if historyFlags.tail > 0 && historyFlags.tail < len(log) {
			log = log[len(log)-historyFlags.tail:]
		}

		if len(log) == 0 {
			fmt.Printf("no tasks for agent %s\n", a.ID)
			return nil
		}

		for _, t := range log {
			status := "ok"
			if t.IsError {
				status = "ERR"
			}
			fmt.Printf("[%s] [%s] %s\n", t.CompletedAt.Format("2006-01-02 15:04:05"), status, t.Prompt)
			if t.Response != "" && t.Prompt != "__spawn__" {
				// Print first 120 chars of response as preview.
				preview := t.Response
				if len(preview) > 120 {
					preview = preview[:120] + "…"
				}
				fmt.Printf("         → %s\n", preview)
			}
		}
		return nil
	},
}

func init() {
	historyCmd.Flags().IntVar(&historyFlags.tail, "tail", 0, "show last N entries (0 = all)")
	rootCmd.AddCommand(historyCmd)
}
