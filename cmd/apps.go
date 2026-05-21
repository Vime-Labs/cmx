package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Gerencia aplicações",
}

// ── list ────────────────────────────────────────────────────────────────────

var appsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todas as aplicações",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		apps, err := client.ListApps()
		if err != nil {
			return err
		}
		if len(apps) == 0 {
			fmt.Println("Nenhuma aplicação encontrada.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "UUID\tNOME\tSTATUS\tREPO\tBRANCH")
		for _, a := range apps {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				a.UUID, a.Name, a.Status, a.Repository, a.Branch)
		}
		return w.Flush()
	},
}

// ── get ─────────────────────────────────────────────────────────────────────

var appsGetCmd = &cobra.Command{
	Use:   "get <uuid|nome>",
	Short: "Exibe detalhes de uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		app, err := client.GetApp(args[0])
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "UUID:\t%s\n", app.UUID)
		fmt.Fprintf(w, "Nome:\t%s\n", app.Name)
		fmt.Fprintf(w, "Status:\t%s\n", app.Status)
		fmt.Fprintf(w, "Repo:\t%s\n", app.Repository)
		fmt.Fprintf(w, "Branch:\t%s\n", app.Branch)
		fmt.Fprintf(w, "Build pack:\t%s\n", app.BuildPack)
		fmt.Fprintf(w, "Domínios:\t%s\n", app.Domains)
		fmt.Fprintf(w, "Criado em:\t%s\n", app.CreatedAt)
		fmt.Fprintf(w, "Atualizado:\t%s\n", app.UpdatedAt)
		return w.Flush()
	},
}

// ── logs ────────────────────────────────────────────────────────────────────

var logsLines int

var appsLogsCmd = &cobra.Command{
	Use:   "logs <uuid|nome>",
	Short: "Exibe logs de uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		logs, err := client.AppLogs(args[0], logsLines)
		if err != nil {
			return err
		}
		fmt.Print(logs)
		return nil
	},
}

// ── deploy ───────────────────────────────────────────────────────────────────

var deployForce bool

var appsDeployCmd = &cobra.Command{
	Use:   "deploy <uuid|nome>",
	Short: "Dispara um deploy",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		resp, err := client.DeployApp(args[0], deployForce)
		if err != nil {
			return err
		}
		for _, d := range resp.Deployments {
			fmt.Printf("Deploy enfileirado: %s\n", d.DeploymentUUID)
		}
		return nil
	},
}

// ── start / stop / restart ──────────────────────────────────────────────────

var appsStartCmd = &cobra.Command{
	Use:   "start <uuid|nome>",
	Short: "Inicia uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE:  appActionCmd("start"),
}

var appsStopCmd = &cobra.Command{
	Use:   "stop <uuid|nome>",
	Short: "Para uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE:  appActionCmd("stop"),
}

var appsRestartCmd = &cobra.Command{
	Use:   "restart <uuid|nome>",
	Short: "Reinicia uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE:  appActionCmd("restart"),
}

func appActionCmd(action string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		var msg string
		var err error
		switch action {
		case "start":
			msg, err = client.StartApp(args[0])
		case "stop":
			msg, err = client.StopApp(args[0])
		case "restart":
			msg, err = client.RestartApp(args[0])
		}
		if err != nil {
			return err
		}
		fmt.Println(msg)
		return nil
	}
}

func init() {
	appsLogsCmd.Flags().IntVarP(&logsLines, "lines", "n", 100, "Número de linhas")
	appsDeployCmd.Flags().BoolVarP(&deployForce, "force", "f", false, "Força rebuild sem cache")

	appsCmd.AddCommand(appsListCmd, appsGetCmd, appsLogsCmd, appsDeployCmd,
		appsStartCmd, appsStopCmd, appsRestartCmd)
	rootCmd.AddCommand(appsCmd)
}
