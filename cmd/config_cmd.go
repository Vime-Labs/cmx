package cmd

import (
	"fmt"
	"strings"

	"github.com/Vime-Labs/cmx/internal/config"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Gerencia a configuração do CMX",
	Long:  "Exibe ou altera a URL e o token do Coolify salvos na configuração.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Exibe a configuração atual",
	Long:  "Mostra a URL e o token configurados. Se uma chave (url ou token) for informada, exibe apenas o valor dela.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Altera um valor da configuração",
	Long:  "Define a URL ou o token do Coolify. Chaves válidas: url, token.",
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		ui.Fail(fmt.Sprintf("Erro ao carregar config: %v", err))
		return nil
	}

	// No config and no env vars — nothing to show
	if cfg.URL == "" && cfg.Token == "" {
		fmt.Println("Configuração não encontrada. Execute `cmx configure`")
		return nil
	}

	if len(args) == 1 {
		key := strings.ToLower(args[0])
		switch key {
		case "url":
			fmt.Println(cfg.URL)
		case "token":
			fmt.Println(cfg.Token)
		default:
			return fmt.Errorf("chave inválida: %s (valores válidos: url, token)", key)
		}
		return nil
	}

	// Show both in a simple key-value layout
	fmt.Printf("url     %s\n", cfg.URL)
	fmt.Printf("token   %s\n", maskToken(cfg.Token))
	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := strings.ToLower(args[0])
	value := args[1]

	if key != "url" && key != "token" {
		return fmt.Errorf("chave inválida: %s (valores válidos: url, token)", key)
	}

	if value == "" {
		return fmt.Errorf("o valor não pode ser vazio")
	}

	// Load existing config (starts empty if none exists)
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}

	switch key {
	case "url":
		cfg.URL = value
	case "token":
		cfg.Token = value
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("salvando config: %w", err)
	}

	ui.Success(fmt.Sprintf("%s atualizado para %s", key, value))
	return nil
}

// maskToken oculta parte do token para exibição segura.
// Exibe os 3 primeiros e 3 últimos caracteres, ex: tok***ken.
func maskToken(token string) string {
	if len(token) <= 6 {
		return strings.Repeat("*", len(token))
	}
	return token[:3] + strings.Repeat("*", len(token)-6) + token[len(token)-3:]
}
