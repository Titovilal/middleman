package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

type ghRelease struct {
	TagName string `json:"tag_name"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update mdm to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Current version: %s\n", Version)

		resp, err := http.Get("https://api.github.com/repos/Titovilal/middleman/releases/latest")
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("GitHub API returned %d", resp.StatusCode)
		}

		var release ghRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return fmt.Errorf("failed to parse release info: %w", err)
		}

		latest := strings.TrimPrefix(release.TagName, "v")
		current := strings.TrimPrefix(Version, "v")
		if latest == current {
			fmt.Println("Already up to date.")
			return nil
		}

		fmt.Printf("New version available: %s\n", latest)

		goos := runtime.GOOS
		goarch := runtime.GOARCH
		asset := fmt.Sprintf("mdm-%s-%s", goos, goarch)
		url := fmt.Sprintf("https://github.com/Titovilal/middleman/releases/download/%s/%s", latest, asset)

		// Download to temp file
		tmp := filepath.Join(os.TempDir(), "mdm-update")
		fmt.Printf("Downloading %s...\n", latest)

		curlCmd := exec.Command("curl", "-sL", url, "-o", tmp)
		curlCmd.Stderr = os.Stderr
		if err := curlCmd.Run(); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		if err := os.Chmod(tmp, 0o755); err != nil {
			return fmt.Errorf("chmod failed: %w", err)
		}

		// Find where the current binary lives
		currentBin, err := os.Executable()
		if err != nil {
			return fmt.Errorf("could not find current binary: %w", err)
		}
		currentBin, _ = filepath.EvalSymlinks(currentBin)

		// Try to replace directly, fall back to sudo
		fmt.Printf("Installing to %s...\n", currentBin)
		if err := os.Rename(tmp, currentBin); err != nil {
			// Rename failed (cross-device or permissions), try sudo mv
			if strings.Contains(currentBin, "/usr/") {
				mvCmd := exec.Command("sudo", "mv", tmp, currentBin)
				mvCmd.Stderr = os.Stderr
				if err := mvCmd.Run(); err != nil {
					return fmt.Errorf("install failed (try running with sudo): %w", err)
				}
			} else {
				return fmt.Errorf("install failed: %w", err)
			}
		}

		fmt.Printf("Updated to %s.\n", latest)
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.TrimPrefix(Version, "v"))
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)
}
