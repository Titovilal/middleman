package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var runFlags struct {
	connector string
}

var runCmd = &cobra.Command{
	Use:   "run [request]",
	Short: "Run the Middleman with a user request",
	Long: fmt.Sprintf(`Launches an AI CLI as a Middleman agent. The Middleman receives
the guide from .mdm/guides/the_middleman.md, is instructed to create
subagents in the background without asking, and processes your request.

Supported connectors: %s`, strings.Join(connectorNames(), ", ")),
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		wd := workDir

		connName := runFlags.connector
		if connName == "" {
			cfg := loadConfig(filepath.Join(wd, ".mdm"))
			connName = cfg.DefaultCLI
		}
		conn, ok := connectors[connName]
		if !ok {
			return fmt.Errorf("unknown connector %q (available: %s)", connName, strings.Join(connectorNames(), ", "))
		}

		middlemanGuide, err := os.ReadFile(filepath.Join(wd, ".mdm", "guides", "the_middleman.md"))
		if err != nil {
			return fmt.Errorf("read the_middleman.md: %w (run 'mdm init' first)", err)
		}

		userRequest := strings.Join(args, " ")

		prompt := fmt.Sprintf(`You are the Middleman Agent. Follow the guide below strictly.
Create subagents in the background. Do not ask for confirmation, just act.

## Middleman Guide
%s

## User Request
%s`, string(middlemanGuide), userRequest)

		fmt.Fprintf(os.Stderr, "Running middleman with %s...\n", conn.Name)

		result, err := conn.Run(wd, prompt)
		if err != nil {
			return err
		}

		fmt.Println(result)
		return nil
	},
}

func init() {
	runCmd.Flags().StringVarP(&runFlags.connector, "connector", "c", "", "AI CLI to use (default: from config)")
	rootCmd.AddCommand(runCmd)
}
