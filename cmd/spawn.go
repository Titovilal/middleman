package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var spawnFlags struct {
	briefing  string
	connector string
}

var spawnCmd = &cobra.Command{
	Use:   "spawn <name>",
	Short: "Create a new agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		connName := spawnFlags.connector
		if connName == "" {
			connName = cfg.DefaultConnector
		}
		briefing := spawnFlags.briefing

		a, err := orch.Spawn(context.Background(), id, briefing, connName)
		if err != nil {
			return err
		}

		fmt.Printf("agent spawned\n")
		fmt.Printf("  id:        %s\n", a.ID)
		fmt.Printf("  connector: %s\n", a.ConnectorName)
		fmt.Printf("  session:   %s\n", a.SessionID)
		fmt.Printf("  status:    %s\n", a.Status)
		if len(a.Checkpoints) > 0 {
			fmt.Printf("  checkpoint: %s\n", a.Checkpoints[0].Label)
		}
		return nil
	},
}

func init() {
	spawnCmd.Flags().StringVarP(&spawnFlags.briefing, "briefing", "b", "", "initial context for the agent")
	spawnCmd.Flags().StringVar(&spawnFlags.connector, "connector", "", "connector to use (overrides default)")
	rootCmd.AddCommand(spawnCmd)
}
