package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// connector defines how to call an AI CLI.
type connector struct {
	Name string
	// Open launches an interactive session with the given prompt.
	Open func(workDir, prompt string) error
	// Run executes the prompt synchronously and returns the text response.
	Run func(workDir, prompt string) (string, error)
}

var connectors = map[string]connector{
	"claude":   {Name: "claude", Open: openClaude, Run: runClaude},
	"copilot":  {Name: "copilot", Open: openCopilot, Run: runCopilot},
	"gemini":   {Name: "gemini", Open: openGemini, Run: runGemini},
	"codex":    {Name: "codex", Open: openCodex, Run: runCodex},
	"opencode": {Name: "opencode", Open: openOpenCode, Run: runOpenCode},
}

// runInteractive launches a CLI command attached to the terminal (stdin/stdout/stderr).
func runInteractive(workDir string, name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Dir = workDir
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// runSync launches a CLI command and captures its stdout output.
func runSync(workDir string, name string, args ...string) (string, error) {
	c := exec.Command(name, args...)
	c.Dir = workDir
	c.Stderr = os.Stderr
	out, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("%s failed: %w", name, err)
	}
	return string(out), nil
}

// --- Claude ---

func openClaude(workDir, prompt string) error {
	return runInteractive(workDir, "claude", "--dangerously-skip-permissions", prompt)
}

type claudeOutput struct {
	Result  string `json:"result"`
	IsError bool   `json:"is_error"`
}

func runClaude(workDir, prompt string) (string, error) {
	c := exec.Command("claude", "--print", "--output-format", "json", "--dangerously-skip-permissions", prompt)
	c.Dir = workDir
	c.Stderr = os.Stderr

	var stdout bytes.Buffer
	c.Stdout = &stdout

	if err := c.Run(); err != nil {
		if stdout.Len() == 0 {
			return "", fmt.Errorf("claude failed: %w", err)
		}
	}

	var out claudeOutput
	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil {
		return "", fmt.Errorf("failed to parse claude output: %v\nraw: %s", err, stdout.String())
	}
	if out.IsError {
		return "", fmt.Errorf("claude returned error: %s", out.Result)
	}
	return out.Result, nil
}

// --- Copilot (GitHub Copilot CLI) ---

func openCopilot(workDir, prompt string) error {
	return runInteractive(workDir, "copilot", "--prompt", prompt)
}

func runCopilot(workDir, prompt string) (string, error) {
	return runSync(workDir, "copilot", "--prompt", prompt, "--silent")
}

// --- Gemini ---

func openGemini(workDir, prompt string) error {
	return runInteractive(workDir, "gemini", prompt)
}

func runGemini(workDir, prompt string) (string, error) {
	return runSync(workDir, "gemini", "--noinput", prompt)
}

// --- Codex ---

func openCodex(workDir, prompt string) error {
	return runInteractive(workDir, "codex", "--approval-mode", "full-auto", prompt)
}

func runCodex(workDir, prompt string) (string, error) {
	return runSync(workDir, "codex", "--approval-mode", "full-auto", "--quiet", prompt)
}

// --- OpenCode ---

func openOpenCode(workDir, prompt string) error {
	return runInteractive(workDir, "opencode", prompt)
}

func runOpenCode(workDir, prompt string) (string, error) {
	return runSync(workDir, "opencode", "--non-interactive", prompt)
}
