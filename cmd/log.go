package cmd

import (
	"fmt"

	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var (
	logN         int
	logCmdFilter string
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Exibe o histórico de atividade da CLI",
	Long: `Mostra um histórico local de tudo que foi feito usando o cmx.

Exemplos:
  cmx log                     # últimas 20 ações
  cmx log -n 50               # últimas 50 ações
  cmx log --cmd deploy        # actions apenas de deploy
  cmx log --cmd apps          # tudo relacionado a apps
  cmx log --clear             # apaga todo o histórico`,
	RunE: runLog,
}

var logClear bool

func runLog(cmd *cobra.Command, args []string) error {
	if logClear {
		if !ui.Confirm("Apagar todo o histórico de atividade?") {
			fmt.Println("Cancelado.")
			return nil
		}
		if err := logger.Clear(); err != nil {
			return fmt.Errorf("erro ao limpar histórico: %w", err)
		}
		ui.Success("Histórico apagado")
		return nil
	}

	entries, err := logger.List(logN, logCmdFilter)
	if err != nil {
		return fmt.Errorf("erro ao ler histórico: %w", err)
	}

	if len(entries) == 0 {
		ui.Info("Nenhuma atividade registrada ainda.")
		return nil
	}

	// Estatísticas rápidas
	stats, _ := logger.ComputeStats()
	fmt.Println()
	fmt.Printf("  %s  %d total  ·  %d sucesso  ·  %d erro\n",
		ui.Bold("Histórico de atividade"),
		stats.Total, stats.Success, stats.Errors,
	)
	fmt.Println()

	// Exibe agrupado por dia
	groups := logger.GroupByDay(entries)
	for _, day := range logger.SortedDays(groups) {
		fmt.Printf("  %s\n", ui.Bold(day))
		for _, e := range groups[day] {
			statusColor := ui.Green("✓")
			if e.Status != "success" {
				statusColor = ui.Red("✗")
			}

			time := e.Timestamp[11:19] // HH:MM:SS
			dur := ""
			if e.Duration != "" {
				dur = ui.Gray(" (" + e.Duration + ")")
			}

			detail := ""
			if e.Detail != "" {
				detail = " — " + e.Detail
			}

			// Formata:  09:30  ✓  deploy meu-app — 3 deployment(s) enfileirado(s) (2.3s)
			fmt.Printf("    %s  %s  %s %s%s%s\n",
				ui.Gray(time),
				statusColor,
				ui.Cyan(e.Command),
				e.ResourceName,
				detail,
				dur,
			)
		}
	}
	fmt.Println()
	return nil
}

func init() {
	logCmd.Flags().IntVarP(&logN, "n", "n", 20, "Número de entradas a exibir (0 = todas)")
	logCmd.Flags().StringVar(&logCmdFilter, "cmd", "", "Filtrar por comando (ex: deploy, apps, dbs)")
	logCmd.Flags().BoolVar(&logClear, "clear", false, "Apagar todo o histórico")
	rootCmd.AddCommand(logCmd)
}
