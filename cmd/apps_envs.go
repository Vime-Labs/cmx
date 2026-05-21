package cmd

import (
	"fmt"
	"os"

	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/Vime-Labs/cmx/internal/validate"
	"github.com/spf13/cobra"
)

var appsEnvsCmd = &cobra.Command{
	Use:   "envs",
	Short: "Gerencia variáveis de ambiente de uma aplicação",
}

// ── list ────────────────────────────────────────────────────────────────────

var appsEnvsListCmd = &cobra.Command{
	Use:   "list <uuid|nome>",
	Short: "Lista variáveis de ambiente",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando variáveis")
		envs, err := client.ListAppEnvs(args[0])
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%d variável(is) encontrada(s)", len(envs)))

		if len(envs) == 0 {
			return nil
		}
		fmt.Println()
		t := ui.NewTable("CHAVE", "VALOR", "PREVIEW")
		for _, e := range envs {
			preview := ui.Gray("não")
			if e.IsPreview {
				preview = ui.Cyan("sim")
			}
			// mascarar valores que parecem segredos
			value := maskSecret(e.Key, e.Value)
			t.AddRow(ui.Bold(e.Key), value, preview)
		}
		t.Render()
		return nil
	},
}

// ── set ─────────────────────────────────────────────────────────────────────

var appsEnvsSetCmd = &cobra.Command{
	Use:   "set <uuid|nome> KEY=VALUE",
	Short: "Cria ou atualiza uma variável de ambiente",
	Long: `Cria ou atualiza uma variável de ambiente.

Exemplos:
  cmx apps envs set meu-app DATABASE_URL=postgres://...
  cmx apps envs set meu-app NODE_ENV production`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID := args[0]
		key, value, err := validate.KeyValue(args[1:])
		if err != nil {
			return err
		}

		client := mustClient()
		spin := ui.NewSpinner(fmt.Sprintf("Definindo %s", key))
		if err := client.SetAppEnv(appID, key, value); err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%s definida", key))
		return nil
	},
}

// ── delete ───────────────────────────────────────────────────────────────────

var appsEnvsDeleteCmd = &cobra.Command{
	Use:   "delete <uuid|nome> KEY",
	Short: "Remove uma variável de ambiente",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		appID, key := args[0], args[1]

		client := mustClient()
		spin := ui.NewSpinner(fmt.Sprintf("Removendo %s", key))
		if err := client.DeleteAppEnvByKey(appID, key); err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%s removida", key))
		return nil
	},
}

// maskSecret oculta parcialmente valores de chaves que parecem segredos.
func maskSecret(key, value string) string {
	sensitive := []string{"KEY", "SECRET", "TOKEN", "PASSWORD", "PASS", "PWD", "PRIVATE", "CREDENTIAL"}
	keyUp := toUpper(key)
	for _, s := range sensitive {
		if contains(keyUp, s) {
			if len(value) <= 4 {
				return ui.Gray("****")
			}
			return value[:4] + ui.Gray("****")
		}
	}
	return value
}

func toUpper(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			result[i] = c - 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func init() {
	_ = os.Stderr // garante que os é usado se necessário
	appsEnvsCmd.AddCommand(appsEnvsListCmd, appsEnvsSetCmd, appsEnvsDeleteCmd)
	appsCmd.AddCommand(appsEnvsCmd)
}
