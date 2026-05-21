package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var dbsCmd = &cobra.Command{
	Use:   "dbs",
	Short: "Gerencia bancos de dados",
}

// ── list ────────────────────────────────────────────────────────────────────

var dbsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos os bancos de dados",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando bancos de dados")
		dbs, err := client.ListDBs()
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%d banco(s) encontrado(s)", len(dbs)))

		if len(dbs) == 0 {
			return nil
		}
		fmt.Println()
		t := ui.NewTable("UUID", "NOME", "TIPO", "IMAGEM", "STATUS", "PÚBLICO")
		for _, d := range dbs {
			public := ui.Gray("—")
			if d.IsPublic {
				public = strconv.Itoa(d.PublicPort)
			}
			t.AddRow(
				ui.Gray(ui.ShortID(d.UUID)),
				ui.Truncate(d.Name, 25),
				ui.Cyan(d.DisplayType()),
				ui.Gray(d.Image),
				ui.StatusColor(d.Status),
				public,
			)
		}
		t.Render()
		return nil
	},
}

// ── get ─────────────────────────────────────────────────────────────────────

var dbsGetCmd = &cobra.Command{
	Use:   "get <uuid|nome>",
	Short: "Exibe detalhes de um banco de dados",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando")
		db, err := client.GetDB(args[0])
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(db.Name)
		fmt.Println()

		public := ui.Gray("não")
		if db.IsPublic {
			public = fmt.Sprintf("sim (porta %d)", db.PublicPort)
		}
		kv := ui.NewTable("", "")
		kv.AddRow(ui.Bold("UUID:"), db.UUID)
		kv.AddRow(ui.Bold("Nome:"), db.Name)
		kv.AddRow(ui.Bold("Tipo:"), ui.Cyan(db.DisplayType()))
		kv.AddRow(ui.Bold("Imagem:"), db.Image)
		kv.AddRow(ui.Bold("Status:"), ui.StatusColor(db.Status))
		kv.AddRow(ui.Bold("Público:"), public)
		kv.AddRow(ui.Bold("Criado em:"), db.CreatedAt)
		kv.AddRow(ui.Bold("Atualizado:"), db.UpdatedAt)
		kv.Render()
		return nil
	},
}

// ── start / stop / restart ──────────────────────────────────────────────────

var dbsStartCmd = &cobra.Command{
	Use:   "start <uuid|nome>",
	Short: "Inicia um banco de dados",
	Args:  cobra.ExactArgs(1),
	RunE:  dbActionCmd("start", "Iniciando"),
}

var dbsStopCmd = &cobra.Command{
	Use:   "stop <uuid|nome>",
	Short: "Para um banco de dados",
	Args:  cobra.ExactArgs(1),
	RunE:  dbActionCmd("stop", "Parando"),
}

var dbsRestartCmd = &cobra.Command{
	Use:   "restart <uuid|nome>",
	Short: "Reinicia um banco de dados",
	Args:  cobra.ExactArgs(1),
	RunE:  dbActionCmd("restart", "Reiniciando"),
}

var dbActionMap = map[string]logger.Action{
	"start":   logger.ActionDBStart,
	"stop":    logger.ActionDBStop,
	"restart": logger.ActionDBRestart,
}

func dbActionCmd(action, spinMsg string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("%s %q", spinMsg, args[0]))
		var msg string
		var err error
		switch action {
		case "start":
			msg, err = client.StartDB(args[0])
		case "stop":
			msg, err = client.StopDB(args[0])
		case "restart":
			msg, err = client.RestartDB(args[0])
		}
		if err != nil {
			spin.Fail("falhou")
			logger.Log(dbActionMap[action], logger.ResourceDB, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop(msg)
		logger.Log(dbActionMap[action], logger.ResourceDB, args[0], msg, "success", time.Since(start))
		return nil
	}
}

func init() {
	dbsCmd.AddCommand(dbsListCmd, dbsGetCmd, dbsStartCmd, dbsStopCmd, dbsRestartCmd)
	rootCmd.AddCommand(dbsCmd)
}
