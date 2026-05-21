package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

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
		dbs, err := client.ListDBs()
		if err != nil {
			return err
		}
		if len(dbs) == 0 {
			fmt.Println("Nenhum banco de dados encontrado.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "UUID\tNOME\tTIPO\tIMAGEM\tSTATUS\tPÚBLICO")
		for _, d := range dbs {
			public := "-"
			if d.IsPublic {
				public = strconv.Itoa(d.PublicPort)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				d.UUID, d.Name, d.Type, d.Image, d.Status, public)
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
		db, err := client.GetDB(args[0])
		if err != nil {
			return err
		}
		public := "não"
		if db.IsPublic {
			public = fmt.Sprintf("sim (porta %d)", db.PublicPort)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "UUID:\t%s\n", db.UUID)
		fmt.Fprintf(w, "Nome:\t%s\n", db.Name)
		fmt.Fprintf(w, "Tipo:\t%s\n", db.Type)
		fmt.Fprintf(w, "Imagem:\t%s\n", db.Image)
		fmt.Fprintf(w, "Status:\t%s\n", db.Status)
		fmt.Fprintf(w, "Público:\t%s\n", public)
		fmt.Fprintf(w, "Criado em:\t%s\n", db.CreatedAt)
		fmt.Fprintf(w, "Atualizado:\t%s\n", db.UpdatedAt)
		return w.Flush()
	},
}

// ── start / stop / restart ──────────────────────────────────────────────────

var dbsStartCmd = &cobra.Command{
	Use:   "start <uuid|nome>",
	Short: "Inicia um banco de dados",
	Args:  cobra.ExactArgs(1),
	RunE:  dbActionCmd("start"),
}

var dbsStopCmd = &cobra.Command{
	Use:   "stop <uuid|nome>",
	Short: "Para um banco de dados",
	Args:  cobra.ExactArgs(1),
	RunE:  dbActionCmd("stop"),
}

var dbsRestartCmd = &cobra.Command{
	Use:   "restart <uuid|nome>",
	Short: "Reinicia um banco de dados",
	Args:  cobra.ExactArgs(1),
	RunE:  dbActionCmd("restart"),
}

func dbActionCmd(action string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := mustClient()
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
			return err
		}
		fmt.Println(msg)
		return nil
	}
}

func init() {
	dbsCmd.AddCommand(dbsListCmd, dbsGetCmd, dbsStartCmd, dbsStopCmd, dbsRestartCmd)
	rootCmd.AddCommand(dbsCmd)
}
