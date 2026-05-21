package cmd

import (
	"fmt"
	"os"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/config"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "cmx",
	Short:        "Vime Labs — Cloud Management CLI",
	Long:         "CMX gerencia aplicações e bancos de dados nos servidores Vime.",
	SilenceUsage: true, // não imprime usage em erros de runtime
}

func Execute() {
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		ui.Fail(err.Error())
		os.Exit(1)
	}
}

// mustClient carrega config e retorna um cliente pronto.
// Encerra o processo com mensagem amigável se config inválida.
func mustClient() *api.Client {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao carregar config: %v\n", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	return api.NewClient(cfg.URL, cfg.Token)
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
