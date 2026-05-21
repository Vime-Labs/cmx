package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version é injetado em build time via -ldflags.
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão do CMX",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cmx %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
