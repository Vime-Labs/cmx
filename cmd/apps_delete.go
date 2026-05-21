package cmd

import (
	"fmt"
	"time"

	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var appsDeleteCmd = &cobra.Command{
	Use:               "delete <uuid|nome>",
	Short:             "Remove uma aplicação",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeApps,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()

		// Confirmação antes de deletar
		if !ui.Confirm(fmt.Sprintf("Remover aplicação %q?", args[0])) {
			fmt.Println("Cancelado.")
			return nil
		}

		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("Removendo %q", args[0]))
		if err := client.DeleteApp(args[0]); err != nil {
			spin.Fail("falhou")
			logger.Log(logger.ActionAppDelete, logger.ResourceApp, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop("Aplicação removida")
		logger.Log(logger.ActionAppDelete, logger.ResourceApp, args[0], "", "success", time.Since(start))
		return nil
	},
}

func init() {
	appsCmd.AddCommand(appsDeleteCmd)
}
