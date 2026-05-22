package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/config"
	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/spf13/cobra"
)

var (
	configureURL   string
	configureToken string
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configura a URL do Coolify e o token de acesso",
	RunE:  runConfigure,
}

func runConfigure(cmd *cobra.Command, args []string) error {
	url := configureURL
	token := configureToken

	if url == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("URL do Coolify (ex: http://192.168.1.10:8000): ")
		url, _ = reader.ReadString('\n')
		url = strings.TrimSpace(url)
	}
	if token == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Token de API: ")
		token, _ = reader.ReadString('\n')
		token = strings.TrimSpace(token)
	}

	if url == "" || token == "" {
		return fmt.Errorf("URL e token são obrigatórios")
	}

	fmt.Print("Verificando conexão... ")
	client := api.NewClient(url, token)
	if err := client.Ping(); err != nil {
		fmt.Println("falhou")
		return fmt.Errorf("não foi possível conectar: %w", err)
	}
	fmt.Println("OK")

	cfg := &config.Config{URL: url, Token: token}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("salvando config: %w", err)
	}

	logger.Log(logger.ActionConfigure, logger.ResourceCfg, url, "", "success")
	fmt.Println("Configuração salva em ~/.cmx/config.yaml")
	return nil
}

func init() {
	configureCmd.Flags().StringVar(&configureURL, "url", "", "URL do servidor Coolify")
	configureCmd.Flags().StringVar(&configureToken, "token", "", "Token de API do Coolify")
}
