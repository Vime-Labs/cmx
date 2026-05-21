package cmd

import (
	"fmt"
	"time"

	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Gerencia domínios de aplicações",
}

// ── list ────────────────────────────────────────────────────────────────────

var domainListCmd = &cobra.Command{
	Use:               "list <uuid|nome>",
	Short:             "Lista domínios de uma aplicação",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeApps,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner(fmt.Sprintf("Buscando domínios de %q", args[0]))
		domains, err := client.ListDomains(args[0])
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%d domínio(s) encontrado(s)", len(domains)))

		if renderJSON(domains) {
			return nil
		}

		if len(domains) == 0 {
			return nil
		}
		fmt.Println()
		t := ui.NewTable("UUID", "DOMÍNIO", "TIPO", "CRIADO EM")
		for _, d := range domains {
			t.AddRow(
				ui.Gray(ui.ShortID(d.UUID)),
				d.Domain,
				ui.Cyan(d.Type),
				ui.Gray(d.CreatedAt),
			)
		}
		t.Render()
		return nil
	},
}

// ── add ─────────────────────────────────────────────────────────────────────

var domainAddCmd = &cobra.Command{
	Use:               "add <uuid|nome> <domínio>",
	Short:             "Adiciona um domínio a uma aplicação",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeApps,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("Adicionando domínio %q a %q", args[1], args[0]))
		d, err := client.AddDomain(args[0], args[1])
		if err != nil {
			spin.Fail("falhou")
			logger.Log(logger.ActionDomainAdd, logger.ResourceApp, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop("Domínio adicionado")
		logger.Log(logger.ActionDomainAdd, logger.ResourceApp, args[0],
			fmt.Sprintf("domínio: %s (UUID: %s)", d.Domain, d.UUID), "success", time.Since(start))
		fmt.Println()
		kv := ui.NewTable("", "")
		kv.AddRow(ui.Bold("UUID:"), d.UUID)
		kv.AddRow(ui.Bold("Domínio:"), d.Domain)
		kv.AddRow(ui.Bold("Tipo:"), ui.Cyan(d.Type))
		kv.AddRow(ui.Bold("Criado em:"), d.CreatedAt)
		kv.Render()
		return nil
	},
}

// ── remove ──────────────────────────────────────────────────────────────────

var domainRemoveCmd = &cobra.Command{
	Use:               "remove <uuid|nome> <domain-uuid>",
	Short:             "Remove um domínio de uma aplicação",
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeApps,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("Removendo domínio %q de %q", args[1], args[0]))
		err := client.RemoveDomain(args[0], args[1])
		if err != nil {
			spin.Fail("falhou")
			logger.Log(logger.ActionDomainRemove, logger.ResourceApp, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop("Domínio removido")
		logger.Log(logger.ActionDomainRemove, logger.ResourceApp, args[0],
			fmt.Sprintf("domínio UUID: %s", args[1]), "success", time.Since(start))
		return nil
	},
}

func init() {
	domainCmd.AddCommand(domainListCmd, domainAddCmd, domainRemoveCmd)
	rootCmd.AddCommand(domainCmd)
}
