package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Titovilal/middleman/config"
	"github.com/Titovilal/middleman/connector"
	claudeconn "github.com/Titovilal/middleman/connector/claude"
	codexconn "github.com/Titovilal/middleman/connector/codex"
	geminiconn "github.com/Titovilal/middleman/connector/gemini"
	opencodeconn "github.com/Titovilal/middleman/connector/opencode"
	"github.com/Titovilal/middleman/orchestrator"
	"github.com/Titovilal/middleman/store"
	"github.com/spf13/cobra"
)

var defaultsFS embed.FS

// SetDefaultsFS receives the embedded .mdm/ defaults from main.go.
func SetDefaultsFS(fs embed.FS) { defaultsFS = fs }

var (
	cfg   *config.Config
	orch  *orchestrator.Orchestrator
	flags struct {
		workDir   string
		connector string
		global    bool
	}
)

var rootCmd = &cobra.Command{
	Use:   "mdm",
	Short: "The Middleman - AI agent orchestrator",
	Long: `MDM (The Middleman) orchestrates multiple AI agent instances (Claude, Codex, Gemini, OpenCode).
The Middleman manages agent lifecycle, context, checkpoints, and rewinds.`,
	Run: func(cmd *cobra.Command, args []string) {
		printBanner()
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip heavy init for the root command and any command that
		// annotates itself with skip_init (no orchestrator needed).
		if cmd.Name() == "mdm" || cmd.Annotations["skip_init"] == "true" {
			return nil
		}

		cfg = config.Default()
		if flags.workDir != "" {
			cfg.WorkDir = flags.workDir
		}
		if flags.connector != "" {
			cfg.DefaultConnector = flags.connector
		}
		cfg.GlobalMode = flags.global

		var s *store.Store
		var err error
		if cfg.GlobalMode {
			s, err = store.NewGlobal(defaultsFS)
		} else {
			s, err = store.New(cfg.WorkDir, defaultsFS)
		}
		if err != nil {
			return fmt.Errorf("init store: %w", err)
		}

		reg := connector.NewConnectorRegistry()
		reg.Register(claudeconn.New())
		reg.Register(codexconn.New())
		reg.Register(geminiconn.New())
		reg.Register(opencodeconn.New())

		orch = orchestrator.New(s, reg, cfg.WorkDir)

		// If .mdm/ is not populated yet, only sync-docs is allowed.
		guidePath := filepath.Join(cfg.WorkDir, ".mdm", "guides", "how_mdm_works.md")
		if _, err := os.Stat(guidePath); os.IsNotExist(err) && cmd.Name() != "sync-docs" {
			return fmt.Errorf(".mdm/ is not initialized. Run 'mdm sync-docs' first")
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flags.workDir, "workdir", "w", "", "project directory (default: current dir)")
	rootCmd.PersistentFlags().StringVarP(&flags.connector, "connector", "c", "", "default connector (claude|codex|gemini|opencode)")
	rootCmd.PersistentFlags().BoolVarP(&flags.global, "global", "g", false, "use ~/.mdm/ instead of ./.mdm/")
}

func printBanner() {
	const (
		reset  = "\033[0m"
		bold   = "\033[1m"
		dim    = "\033[2m"
		cyan   = "\033[36m"
		blue   = "\033[34m"
		purple = "\033[35m"
		white  = "\033[97m"
		green  = "\033[32m"

		// Spanish flag palette (soft tones via 256-color)
		red    = "\033[38;5;167m" // soft red
		yellow = "\033[38;5;220m" // warm yellow
	)

	fmt.Println()
	fmt.Println(red + bold + "  ███╗   ███╗██████╗ ███╗   ███╗" + reset)
	fmt.Println(red + bold + "  ████╗ ████║██╔══██╗████╗ ████║" + reset)
	fmt.Println(yellow + bold + "  ██╔████╔██║██║  ██║██╔████╔██║" + reset)
	fmt.Println(yellow + bold + "  ██║╚██╔╝██║██║  ██║██║╚██╔╝██║" + reset)
	fmt.Println(red + bold + "  ██║ ╚═╝ ██║██████╔╝██║ ╚═╝ ██║" + reset)
	fmt.Println(red + bold + "  ╚═╝     ╚═╝╚═════╝ ╚═╝     ╚═╝" + reset)
	fmt.Println()
	fmt.Println(white + bold + "  The Middleman" + reset + dim + " — AI agent orchestrator" + reset)
	fmt.Println()
	fmt.Printf(dim+"  version "+reset+green+"%s"+reset+dim+"  go "+reset+green+"%s"+reset+dim+"  %s/%s"+reset+"\n",
		Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Println()
	fmt.Println(dim + "  Connectors:  " + reset + cyan + "claude" + reset + dim + " · " + reset + cyan + "codex" + reset + dim + " · " + reset + cyan + "gemini" + reset + dim + " · " + reset + cyan + "opencode" + reset)
	fmt.Println()
	fmt.Println(red + bold + "  Quick start" + reset + dim + " — pick your flow:" + reset)
	fmt.Println()
	fmt.Println(white + bold + "  1. Launch mode" + reset + dim + " (mdm launches the agent for you)" + reset)
	fmt.Println(dim + "     $ " + reset + white + "mdm sync-docs" + reset + dim + "          initialize .mdm/ in your project" + reset)
	fmt.Println(dim + "     $ " + reset + white + "mdm launch claude" + reset + dim + "      start Claude as the Middleman" + reset)
	fmt.Println()
	fmt.Println(white + bold + "  2. Agent mode" + reset + dim + " (you're already inside an AI CLI)" + reset)
	fmt.Println(dim + "     Tell your agent: " + reset + cyan + "\"Run mdm agent-prompt, read it, and act as the Middleman\"" + reset)
	fmt.Println()
	fmt.Println(dim + "  Run " + reset + white + "mdm --help" + reset + dim + " for all commands." + reset)
	fmt.Println(dim + "  Docs: " + reset + blue + "https://github.com/Titovilal/middleman" + reset)
	fmt.Println()
}
