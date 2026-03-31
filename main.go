package main

import (
	"embed"

	"github.com/Titovilal/context0/cmd"
)

//go:embed defaults
var defaultsFS embed.FS

func main() {
	cmd.SetDefaultsFS(defaultsFS)
	cmd.Execute()
}
