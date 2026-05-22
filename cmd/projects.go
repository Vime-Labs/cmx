package cmd

import (
	"fmt"

	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Gerencia projetos",
}

// ── list ────────────────────────────────────────────────────────────────────

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos os projetos",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando projetos")
		projects, err := client.ListProjects()
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(fmt.Sprintf("%d projeto(s) encontrado(s)", len(projects)))

		if renderJSON(projects) {
			return nil
		}

		if len(projects) == 0 {
			return nil
		}
		fmt.Println()
		t := ui.NewTable("UUID", "NOME", "AMBIENTES")
		for _, p := range projects {
			envs := ""
			for i, e := range p.Environments {
				if i > 0 {
					envs += ", "
				}
				envs += e.Name
			}
			t.AddRow(ui.Gray(ui.ShortID(p.UUID)), p.Name, ui.Gray(envs))
		}
		t.Render()
		return nil
	},
}

// ── get ─────────────────────────────────────────────────────────────────────

var projectsGetCmd = &cobra.Command{
	Use:   "get <uuid>",
	Short: "Exibe detalhes de um projeto",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Buscando")
		p, err := client.GetProject(args[0])
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop(p.Name)

		if renderJSON(p) {
			return nil
		}

		fmt.Println()
		kv := ui.NewTable("", "")
		kv.AddRow(ui.Bold("UUID:"), p.UUID)
		kv.AddRow(ui.Bold("Nome:"), p.Name)
		kv.AddRow(ui.Bold("Descrição:"), p.Description)
		envNames := ""
		for i, e := range p.Environments {
			if i > 0 {
				envNames += ", "
			}
			envNames += fmt.Sprintf("%s (%s)", e.Name, e.UUID)
		}
		kv.AddRow(ui.Bold("Ambientes:"), ui.Gray(envNames))
		kv.Render()
		return nil
	},
}

// ── create ──────────────────────────────────────────────────────────────────

var projectsCreateCmd = &cobra.Command{
	Use:   "create <nome>",
	Short: "Cria um novo projeto",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner(fmt.Sprintf("Criando projeto %q", args[0]))
		p, err := client.CreateProject(args[0])
		if err != nil {
			spin.Fail("falhou")
			return err
		}
		spin.Stop("Projeto criado")
		fmt.Println()
		ui.Info(fmt.Sprintf("UUID: %s", p.UUID))
		for _, e := range p.Environments {
			ui.Info(fmt.Sprintf("Ambiente: %s (%s)", e.Name, e.UUID))
		}
		return nil
	},
}

func init() {
	projectsCmd.AddCommand(projectsListCmd, projectsGetCmd, projectsCreateCmd)
	rootCmd.AddCommand(projectsCmd)
}
