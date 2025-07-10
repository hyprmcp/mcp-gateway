package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	var opts ServeOptions
	cmd := &cobra.Command{
		Use:   "mcp-proxy",
		Short: "A proxy for MCP servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(cmd.Context(), opts)
		},
	}
	BindServeOptions(cmd, &opts)
	return cmd
}
