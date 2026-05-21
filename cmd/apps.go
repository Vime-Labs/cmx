package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/Vime-Labs/cmx/internal/ui"
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
		spin := ui.NewSpinner("Buscando aplicações")
		apps, err := client.ListApps()
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%d aplicação(ões) encontrada(s)", len(apps)))

		if len(apps) == 0 {
			return nil
		}
		fmt.Println()
		t := ui.NewTable("UUID", "NOME", "STATUS", "REPO", "BRANCH")
		for _, a := range apps {
			t.AddRow(
				ui.Gray(ui.ShortID(a.UUID)),
				ui.Truncate(a.Name, 30),
				ui.StatusColor(a.Status),
				ui.Truncate(a.Repository, 35),
				ui.Gray(a.Branch),
			)
		}
		t.Render()
		return nil
	},
}

// ── get ─────────────────────────────────────────────────────────────────────

var appsGetCmd = &cobra.Command{
	Use:   "get <uuid|nome>",
	Short: "Exibe detalhes de uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando")
		app, err := client.GetApp(args[0])
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(app.Name)
		fmt.Println()

		kv := ui.NewTable("", "")
		kv.AddRow(ui.Bold("UUID:"), app.UUID)
		kv.AddRow(ui.Bold("Nome:"), app.Name)
		kv.AddRow(ui.Bold("Status:"), ui.StatusColor(app.Status))
		kv.AddRow(ui.Bold("Repo:"), app.Repository)
		kv.AddRow(ui.Bold("Branch:"), app.Branch)
		kv.AddRow(ui.Bold("Build pack:"), app.BuildPack)
		kv.AddRow(ui.Bold("Domínios:"), app.Domains)
		kv.AddRow(ui.Bold("Criado em:"), app.CreatedAt)
		kv.AddRow(ui.Bold("Atualizado:"), app.UpdatedAt)
		kv.Render()
		return nil
	},
}

// ── logs ────────────────────────────────────────────────────────────────────

var (
	logsLines  int
	logsFollow bool
)

var appsLogsCmd = &cobra.Command{
	Use:   "logs <uuid|nome>",
	Short: "Exibe logs de uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()

		if !logsFollow {
			spin := ui.NewSpinner("Buscando logs")
			logs, err := client.AppLogs(args[0], logsLines)
			if err != nil {
				spin.Fail("falhou")
				return err
			}
			spin.Stop("OK")
			fmt.Println()
			fmt.Print(logs)
			return nil
		}

		// modo follow: poll a cada 2s, imprime apenas linhas novas
		ui.Info(fmt.Sprintf("Seguindo logs de %q (Ctrl+C para sair)", args[0]))
		fmt.Println()
		var seen string
		for {
			logs, err := client.AppLogs(args[0], 200)
			if err != nil {
				ui.Fail(fmt.Sprintf("erro ao buscar logs: %v", err))
			} else {
				newContent := diffLogs(seen, logs)
				if newContent != "" {
					fmt.Print(newContent)
					seen = logs
				}
			}
			time.Sleep(2 * time.Second)
		}
	},
}

// diffLogs retorna as linhas de novo que não estavam em prev.
func diffLogs(prev, curr string) string {
	if prev == "" {
		return curr
	}
	if !strings.HasSuffix(prev, "\n") {
		prev += "\n"
	}
	idx := strings.LastIndex(curr, strings.TrimSuffix(prev[max(0, len(prev)-200):], "\n"))
	if idx < 0 {
		return curr
	}
	tail := curr[idx+len(prev)-1:]
	if tail == "\n" || tail == "" {
		return ""
	}
	return tail
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ── deploy ───────────────────────────────────────────────────────────────────

var appsDeployCmd = &cobra.Command{
	Use:   "deploy <uuid|nome>",
	Short: "Dispara um deploy",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		force, _ := cmd.Flags().GetBool("force")
		spin := ui.NewSpinner(fmt.Sprintf("Disparando deploy para %q", args[0]))
		resp, err := client.DeployApp(args[0], force)
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop("Deploy enfileirado")
		for _, d := range resp.Deployments {
			ui.Info(fmt.Sprintf("deployment: %s", d.DeploymentUUID))
		}
		return nil
	},
}

// ── start / stop / restart ──────────────────────────────────────────────────

var appsStartCmd = &cobra.Command{
	Use:   "start <uuid|nome>",
	Short: "Inicia uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE:  appActionCmd("start", "Iniciando"),
}

var appsStopCmd = &cobra.Command{
	Use:   "stop <uuid|nome>",
	Short: "Para uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE:  appActionCmd("stop", "Parando"),
}

var appsRestartCmd = &cobra.Command{
	Use:   "restart <uuid|nome>",
	Short: "Reinicia uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE:  appActionCmd("restart", "Reiniciando"),
}

func appActionCmd(action, spinMsg string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner(fmt.Sprintf("%s %q", spinMsg, args[0]))
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
			spin.Fail("falhou")
			return err
		}
		spin.Stop(msg)
		return nil
	}
}

func init() {
	appsLogsCmd.Flags().IntVarP(&logsLines, "lines", "n", 100, "Número de linhas")
	appsLogsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Acompanha logs em tempo real (poll 2s)")
	appsDeployCmd.Flags().Bool("force", false, "Força rebuild sem cache")

	appsCmd.AddCommand(appsListCmd, appsGetCmd, appsLogsCmd, appsDeployCmd,
		appsStartCmd, appsStopCmd, appsRestartCmd)
	rootCmd.AddCommand(appsCmd)
}
