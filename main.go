package main

import (
	"os"

	"github.com/jetski-sh/mcp-proxy/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
