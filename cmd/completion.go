package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Gera script de autocomplete para o shell",
	Long: `Gera script de autocomplete para bash, zsh, fish ou PowerShell.

Para usar:

  bash:
    source <(cmx completion bash)
    # ou adicione ao ~/.bashrc

  zsh:
    source <(cmx completion zsh)
    # ou adicione ao ~/.zshrc

  fish:
    cmx completion fish | source

  PowerShell:
    cmx completion powershell | Out-String | Invoke-Expression`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		if err != nil {
			return fmt.Errorf("gerando completion: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
