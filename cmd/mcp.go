package cmd

import (
	"fmt"

	"github.com/Vime-Labs/cmx/internal/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Inicia servidor MCP para agentes de IA",
	Long: `Inicia um servidor MCP (Model Context Protocol) via stdio.

Permite que agentes de IA (opencode, Claude Desktop, Cursor, Copilot, etc.)
usem o CMX como ferramentas estruturadas com parâmetros tipados.

Configure no seu cliente MCP (ex: opencode.json):

  "mcp": {
    "cmx": {
      "type": "local",
      "command": ["cmx", "mcp"],
      "enabled": true
    }
  }

No Claude Desktop (claude_desktop_config.json):

  {
    "mcpServers": {
      "cmx": {
        "command": "cmx",
        "args": ["mcp"]
      }
    }
  }`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		srv := mcp.NewServer(client, Version)
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "CMX MCP server iniciado no modo stdio")
		return srv.Run()
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
