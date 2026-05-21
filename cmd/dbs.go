package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/Vime-Labs/cmx/internal/ui"
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
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, ui.Bold("UUID\tNOME\tTIPO\tIMAGEM\tSTATUS\tPÚBLICO"))
		for _, d := range dbs {
			public := ui.Gray("—")
			if d.IsPublic {
				public = strconv.Itoa(d.PublicPort)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				ui.Gray(d.UUID), d.Name, ui.Cyan(d.Type),
				ui.Gray(d.Image), ui.StatusColor(d.Status), public)
		}
		return w.Flush()
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
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\n", ui.Bold("UUID:"), db.UUID)
		fmt.Fprintf(w, "%s\t%s\n", ui.Bold("Nome:"), db.Name)
		fmt.Fprintf(w, "%s\t%s\n", ui.Bold("Tipo:"), ui.Cyan(db.Type))
		fmt.Fprintf(w, "%s\t%s\n", ui.Bold("Imagem:"), db.Image)
		fmt.Fprintf(w, "%s\t%s\n", ui.Bold("Status:"), ui.StatusColor(db.Status))
		fmt.Fprintf(w, "%s\t%s\n", ui.Bold("Público:"), public)
		fmt.Fprintf(w, "%s\t%s\n", ui.Bold("Criado em:"), db.CreatedAt)
		fmt.Fprintf(w, "%s\t%s\n", ui.Bold("Atualizado:"), db.UpdatedAt)
		return w.Flush()
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

func dbActionCmd(action, spinMsg string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := mustClient()
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
			return err
		}
		spin.Stop(msg)
		return nil
	}
}

func init() {
	dbsCmd.AddCommand(dbsListCmd, dbsGetCmd, dbsStartCmd, dbsStopCmd, dbsRestartCmd)
	rootCmd.AddCommand(dbsCmd)
}
