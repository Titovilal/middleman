package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var launchFlags struct {
	workDir string
}

// launchSpec defines how to start an interactive session for a given connector.
type launchSpec struct {
	bin       string
	buildArgs func(guidePath string) []string
}

var launchSpecs = map[string]launchSpec{
	"claude": {
		bin: "claude",
		buildArgs: func(guidePath string) []string {
			return []string{
				"--dangerously-skip-permissions",
				"You are now the Middleman. Always use the mdm CLI — never create or manage agents directly. Read this and act accordingly: @" + guidePath,
			}
		},
	},
	"gemini": {
		bin: "gemini",
		buildArgs: func(guidePath string) []string {
			data, _ := os.ReadFile(guidePath)
			return []string{"--system-prompt", string(data)}
		},
	},
	"codex": {
		bin: "codex",
		buildArgs: func(guidePath string) []string {
			data, _ := os.ReadFile(guidePath)
			return []string{"--instructions", string(data)}
		},
	},
}

var launchCmd = &cobra.Command{
	Use:         "launch <connector>",
	Short:       "Launch an interactive AI CLI session with the Middleman prompt injected",
	Annotations: map[string]string{"skip_init": "true"},
	Long: `Launches the specified AI CLI tool (claude, gemini, codex) in interactive mode
with the Middleman orchestrator prompt automatically injected. The user gets
a fully interactive session where the AI is already acting as the Middleman.

Examples:
  mdm launch claude
  mdm launch gemini
  mdm launch codex`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		connName := args[0]

		spec, ok := launchSpecs[connName]
		if !ok {
			names := make([]string, 0, len(launchSpecs))
			for n := range launchSpecs {
				names = append(names, n)
			}
			sort.Strings(names)
			return fmt.Errorf("unknown connector %q (available: %s)", connName, strings.Join(names, ", "))
		}

		workDir := launchFlags.workDir
		if workDir == "" {
			workDir, _ = os.Getwd()
		}

		// Check that .mdm/ is initialized.
		guidePath := filepath.Join(workDir, ".mdm", "guides", "how_mdm_works.md")
		if _, err := os.Stat(guidePath); os.IsNotExist(err) {
			return fmt.Errorf(".mdm/ is not initialized. Run 'mdm sync-docs' first")
		}

		// Build the command — pass the guide .md directly.
		cliArgs := spec.buildArgs(guidePath)
		cliPath, err := exec.LookPath(spec.bin)
		if err != nil {
			return fmt.Errorf("%s CLI not found in PATH: %w", spec.bin, err)
		}

		c := exec.Command(cliPath, cliArgs...)
		c.Dir = workDir
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		fmt.Fprintf(os.Stderr, "Launching %s as Middleman orchestrator...\n", connName)

		if err := c.Run(); err != nil {
			// If the CLI exited with a non-zero code, propagate it.
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			return fmt.Errorf("failed to run %s: %w", spec.bin, err)
		}
		return nil
	},
}

func init() {
	launchCmd.Flags().StringVarP(&launchFlags.workDir, "workdir", "w", "", "project directory (default: current dir)")
	rootCmd.AddCommand(launchCmd)
}
