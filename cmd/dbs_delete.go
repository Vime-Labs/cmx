package cmd

import (
	"fmt"
	"time"

	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var dbsDeleteCmd = &cobra.Command{
	Use:               "delete <uuid|nome>",
	Short:             "Remove um banco de dados",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeDBs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()

		if !ui.Confirm(fmt.Sprintf("Remover banco de dados %q?", args[0])) {
			fmt.Println("Cancelado.")
			return nil
		}

		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("Removendo %q", args[0]))
		if err := client.DeleteDB(args[0]); err != nil {
			spin.Fail("falhou")
			logger.Log(logger.ActionDBDelete, logger.ResourceDB, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop("Banco de dados removido")
		logger.Log(logger.ActionDBDelete, logger.ResourceDB, args[0], "", "success", time.Since(start))
		return nil
	},
}

func init() {
	dbsCmd.AddCommand(dbsDeleteCmd)
}
