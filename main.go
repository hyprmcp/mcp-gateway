package main

import (
	"os"

	_ "crypto/sha256"

	"github.com/jetski-sh/mcp-gateway/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
