package config

import (
	"os"
	"path/filepath"
)

// Config holds MDM runtime configuration.
type Config struct {
	// WorkDir is the project directory. Registry is stored in WorkDir/.mdm/.
	WorkDir string

	// DefaultConnector is the connector used when --connector is not specified.
	DefaultConnector string

	// GlobalMode uses ~/.mdm/ instead of WorkDir/.mdm/.
	GlobalMode bool
}

func Default() *Config {
	wd, _ := os.Getwd()
	return &Config{
		WorkDir:          wd,
		DefaultConnector: "claude",
	}
}

func (c *Config) RegistryDir() string {
	if c.GlobalMode {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".mdm")
	}
	return filepath.Join(c.WorkDir, ".mdm")
}
