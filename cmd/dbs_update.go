package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var (
	dbUpdateName        string
	dbUpdateImage       string
	dbUpdatePublic      string
	dbUpdatePublicPort  int
	dbUpdateDescription string
)

var dbsUpdateCmd = &cobra.Command{
	Use:               "update <uuid|nome>",
	Short:             "Atualiza configurações de um banco de dados",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeDBs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()

		req := api.UpdateDBRequest{}
		if cmd.Flags().Changed("name") {
			v := dbUpdateName
			req.Name = &v
		}
		if cmd.Flags().Changed("image") {
			v := dbUpdateImage
			req.Image = &v
		}
		if cmd.Flags().Changed("public") {
			v := dbUpdatePublic == "true" || dbUpdatePublic == "yes" || dbUpdatePublic == "1"
			req.IsPublic = &v
		}
		if cmd.Flags().Changed("public-port") {
			v := dbUpdatePublicPort
			req.PublicPort = &v
		}
		if cmd.Flags().Changed("description") {
			v := dbUpdateDescription
			req.Description = &v
		}

		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("Atualizando %q", args[0]))
		db, err := client.UpdateDB(args[0], req)
		if err != nil {
			spin.Fail("falhou")
			logger.Log(logger.ActionDBUpdate, logger.ResourceDB, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop("Banco de dados atualizado")
		logger.Log(logger.ActionDBUpdate, logger.ResourceDB, args[0], db.Name, "success", time.Since(start))

		fmt.Println()
		kv := ui.NewTable("", "")
		kv.AddRow(ui.Bold("UUID:"), db.UUID)
		kv.AddRow(ui.Bold("Nome:"), db.Name)
		kv.AddRow(ui.Bold("Imagem:"), db.Image)
		kv.AddRow(ui.Bold("Tipo:"), db.DisplayType())
		kv.AddRow(ui.Bold("Público:"), strconv.FormatBool(db.IsPublic))
		kv.AddRow(ui.Bold("Status:"), ui.StatusColor(db.Status))
		kv.Render()
		return nil
	},
}

func init() {
	dbsUpdateCmd.Flags().StringVar(&dbUpdateName, "name", "", "Novo nome")
	dbsUpdateCmd.Flags().StringVar(&dbUpdateImage, "image", "", "Imagem Docker")
	dbsUpdateCmd.Flags().StringVar(&dbUpdatePublic, "public", "", "Tornar público (true/false)")
	dbsUpdateCmd.Flags().IntVar(&dbUpdatePublicPort, "public-port", 0, "Porta pública")
	dbsUpdateCmd.Flags().StringVar(&dbUpdateDescription, "description", "", "Descrição")
	dbsCmd.AddCommand(dbsUpdateCmd)
}
