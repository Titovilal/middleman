package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var syncDocsFlags struct {
	connector string
	workers   int
}

var syncDocsCmd = &cobra.Command{
	Use:   "sync-docs",
	Short: "Create or update .ctx/docs/ using an AI CLI",
	Long: fmt.Sprintf(`Reads the codebase and creates or updates the documentation in .ctx/docs/
following the guide in .ctx/guides/how_to_sync_docs.md and the templates
in .ctx/templates/.

Phase 1: generates project_overview.md (single agent).
Phase 2: generates all other docs in parallel (one agent per doc).

Supported connectors: %s`, strings.Join(connectorNames(), ", ")),
	RunE: func(cmd *cobra.Command, args []string) error {
		wd := workDir

		connName := syncDocsFlags.connector
		if connName == "" {
			cfg := loadConfig(filepath.Join(wd, ".ctx"))
			connName = cfg.DefaultCLI
		}
		conn, ok := connectors[connName]
		if !ok {
			return fmt.Errorf("unknown connector %q (available: %s)", connName, strings.Join(connectorNames(), ", "))
		}

		guide, err := os.ReadFile(filepath.Join(wd, ".ctx", "guides", "how_to_sync_docs.md"))
		if err != nil {
			return fmt.Errorf("read how_to_sync_docs.md: %w (run 'ctx init' first)", err)
		}
		docTemplate, _ := os.ReadFile(filepath.Join(wd, ".ctx", "templates", "doc_template.md"))
		overviewTemplate, _ := os.ReadFile(filepath.Join(wd, ".ctx", "templates", "project_overview_template.md"))

		// ── Phase 1: generate project_overview.md ──

		overviewPrompt := fmt.Sprintf(`You are a documentation agent. Your ONLY job is to read the codebase and create/update .ctx/docs/project_overview.md following the guide and template below.

IMPORTANT: Do NOT create or modify any other doc files. Only project_overview.md.

In the "Documentation available" section, list every doc that SHOULD exist in .ctx/docs/ (besides project_overview.md itself). Each entry must follow this exact format:
- **%s[doc_name].md%s** — [short description]

## Guide
%s

## Project overview template
%s`, "`", "`", string(guide), string(overviewTemplate))

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, cDim+"  ▸ "+cReset+"Phase 1: generating project overview with "+stValue(conn.Name)+stDim("..."))

		result, err := conn.Run(wd, overviewPrompt)
		if err != nil {
			return fmt.Errorf("phase 1 (overview): %w", err)
		}
		fmt.Fprintln(os.Stderr, stOk("  ✓ ")+"Project overview done.")

		// Read the generated overview to extract doc list
		overview, err := os.ReadFile(filepath.Join(wd, ".ctx", "docs", "project_overview.md"))
		if err != nil {
			// If the agent printed it instead of writing it, use the result
			overview = []byte(result)
		}

		docs := parseDocList(string(overview))
		if len(docs) == 0 {
			fmt.Fprintln(os.Stderr, stDim("  ▸ ")+"No additional docs listed in overview. Done.")
			return nil
		}

		// ── Phase 2: generate each doc in parallel ──

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, cDim+"  ▸ "+cReset+fmt.Sprintf("Phase 2: generating %d docs in parallel (workers: %d)", len(docs), syncDocsFlags.workers))

		var mu sync.Mutex
		g := new(errgroup.Group)
		g.SetLimit(syncDocsFlags.workers)

		for _, doc := range docs {
			doc := doc // capture
			g.Go(func() error {
				mu.Lock()
				fmt.Fprintln(os.Stderr, stDim("    ▸ ")+"Starting "+stValue(doc))
				mu.Unlock()

				docPrompt := fmt.Sprintf(`You are a documentation agent. Your ONLY job is to read the codebase and create/update the file .ctx/docs/%s following the guide and template below.

IMPORTANT: Do NOT create or modify any other files. Only .ctx/docs/%s.

Here is the project overview for context:
%s

## Guide
%s

## Doc template
%s`, doc, doc, string(overview), string(guide), string(docTemplate))

				_, err := conn.Run(wd, docPrompt)

				mu.Lock()
				if err != nil {
					fmt.Fprintln(os.Stderr, stErr("    ✗ ")+doc+": "+err.Error())
				} else {
					fmt.Fprintln(os.Stderr, stOk("    ✓ ")+doc)
				}
				mu.Unlock()

				return err
			})
		}

		if err := g.Wait(); err != nil {
			return fmt.Errorf("phase 2: %w", err)
		}

		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, stOk("  ✓ ")+"All docs synced.")
		return nil
	},
}

// parseDocList extracts doc filenames from the "Documentation available" section
// of a project overview. It looks for lines matching: **`doc_name.md`**
var docListRe = regexp.MustCompile("`([a-zA-Z0-9_-]+\\.md)`")

func parseDocList(overview string) []string {
	var docs []string
	inSection := false
	for _, line := range strings.Split(overview, "\n") {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "documentation available") {
			inSection = true
			continue
		}
		// Stop at the next heading or end
		if inSection && strings.HasPrefix(line, "#") {
			break
		}
		if inSection {
			if m := docListRe.FindStringSubmatch(line); m != nil {
				name := m[1]
				if name != "project_overview.md" {
					docs = append(docs, name)
				}
			}
		}
	}
	return docs
}

func connectorNames() []string {
	names := make([]string, 0, len(connectors))
	for k := range connectors {
		names = append(names, k)
	}
	return names
}

func init() {
	syncDocsCmd.Flags().StringVarP(&syncDocsFlags.connector, "connector", "c", "", "AI CLI to use (default: claude)")
	syncDocsCmd.Flags().IntVar(&syncDocsFlags.workers, "workers", 5, "max parallel AI agents")
	rootCmd.AddCommand(syncDocsCmd)
}
