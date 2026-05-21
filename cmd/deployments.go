package cmd

import (
	"fmt"
	"time"

	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var deploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "Acompanha deployments",
}

// ── list ────────────────────────────────────────────────────────────────────

var deploymentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista deployments em andamento",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando deployments")
		deps, err := client.ListDeployments()
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%d deployment(s) ativo(s)", len(deps)))

		if len(deps) == 0 {
			return nil
		}
		fmt.Println()
		t := ui.NewTable("UUID", "APP UUID", "STATUS", "CRIADO EM")
		for _, d := range deps {
			t.AddRow(
				ui.Gray(ui.ShortID(d.UUID)),
				ui.Gray(ui.ShortID(d.ApplicationUUID)),
				ui.StatusColor(d.Status),
				d.CreatedAt,
			)
		}
		t.Render()
		return nil
	},
}

// ── get ─────────────────────────────────────────────────────────────────────

var deploymentsGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Exibe detalhes de um deployment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando deployment")
		d, err := client.GetDeployment(args[0])
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(d.UUID)
		fmt.Println()

		kv := ui.NewTable("", "")
		kv.AddRow(ui.Bold("UUID:"), d.UUID)
		kv.AddRow(ui.Bold("App UUID:"), d.ApplicationUUID)
		kv.AddRow(ui.Bold("Status:"), ui.StatusColor(d.Status))
		kv.AddRow(ui.Bold("Criado em:"), d.CreatedAt)
		kv.AddRow(ui.Bold("Atualizado:"), d.UpdatedAt)
		kv.Render()
		return nil
	},
}

// ── cancel ───────────────────────────────────────────────────────────────────

var deploymentsCancelCmd = &cobra.Command{
	Use:   "cancel <uuid>",
	Short: "Cancela um deployment em andamento",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("Cancelando deployment %s", args[0]))
		if err := client.CancelDeployment(args[0]); err != nil {
			spin.Fail("falhou")
			logger.Log(logger.ActionDeployCancel, logger.ResourceDepl, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop("Deployment cancelado")
		logger.Log(logger.ActionDeployCancel, logger.ResourceDepl, args[0], "", "success", time.Since(start))
		return nil
	},
}

// ── history ──────────────────────────────────────────────────────────────────

var (
	historySkip int
	historyTake int
)

var deploymentsHistoryCmd = &cobra.Command{
	Use:   "history <uuid|nome-da-app>",
	Short: "Histórico de deployments de uma aplicação",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando histórico")
		deps, err := client.ListAppDeployments(args[0], historySkip, historyTake)
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%d deployment(s)", len(deps)))

		if len(deps) == 0 {
			return nil
		}
		fmt.Println()
		t := ui.NewTable("UUID", "STATUS", "CRIADO EM", "ATUALIZADO")
		for _, d := range deps {
			t.AddRow(
				ui.Gray(ui.ShortID(d.UUID)),
				ui.StatusColor(d.Status),
				d.CreatedAt,
				d.UpdatedAt,
			)
		}
		t.Render()
		return nil
	},
}

func init() {
	deploymentsHistoryCmd.Flags().IntVar(&historySkip, "skip", 0, "Pular N deployments")
	deploymentsHistoryCmd.Flags().IntVar(&historyTake, "take", 10, "Exibir N deployments")

	deploymentsCmd.AddCommand(deploymentsListCmd, deploymentsGetCmd,
		deploymentsCancelCmd, deploymentsHistoryCmd)
	rootCmd.AddCommand(deploymentsCmd)
}
