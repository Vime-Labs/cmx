package cmd

import (
	"fmt"

	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Gerencia fontes de código (GitHub Apps)",
}

// ── list ────────────────────────────────────────────────────────────────────

var sourcesListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista GitHub Apps configurados no Coolify",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando GitHub Apps")
		apps, err := client.ListGitHubApps()
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%d GitHub App(s) encontrado(s)", len(apps)))

		if renderJSON(apps) {
			return nil
		}

		if len(apps) == 0 {
			fmt.Println()
			ui.Info("Nenhum GitHub App configurado. Adicione um em: Configuracoes > GitHub Apps")
			return nil
		}

		fmt.Println()
		t := ui.NewTable("UUID", "NOME")
		for _, a := range apps {
			t.AddRow(ui.Gray(ui.ShortID(a.UUID)), a.Name)
		}
		t.Render()
		return nil
	},
}

func init() {
	sourcesCmd.AddCommand(sourcesListCmd)
	rootCmd.AddCommand(sourcesCmd)
}
