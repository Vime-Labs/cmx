package cmd

import (
	"fmt"
	"time"

	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var dbsBackupCmd = &cobra.Command{
	Use:               "backup <uuid|nome>",
	Short:             "Dispara um backup de banco de dados",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeDBs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("Iniciando backup de %q", args[0]))
		result, err := client.BackupDB(args[0])
		if err != nil {
			spin.Fail("falhou")
			logger.Log(logger.ActionDBBackup, logger.ResourceDB, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop("Backup concluído")
		logger.Log(logger.ActionDBBackup, logger.ResourceDB, args[0],
			fmt.Sprintf("UUID: %s", result.UUID), "success", time.Since(start))

		fmt.Println()
		kv := ui.NewTable("", "")
		kv.AddRow(ui.Bold("UUID:"), result.UUID)
		if result.Message != "" {
			kv.AddRow(ui.Bold("Mensagem:"), result.Message)
		}
		if result.Size > 0 {
			kv.AddRow(ui.Bold("Tamanho:"), fmt.Sprintf("%d bytes", result.Size))
		}
		if result.Status != "" {
			kv.AddRow(ui.Bold("Status:"), ui.StatusColor(result.Status))
		}
		kv.Render()
		return nil
	},
}

func init() {
	dbsCmd.AddCommand(dbsBackupCmd)
}
