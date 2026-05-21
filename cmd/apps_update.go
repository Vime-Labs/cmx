package cmd

import (
	"fmt"
	"time"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var (
	appUpdateName       string
	appUpdateBuildPack  string
	appUpdateBranch     string
	appUpdateRepository string
	appUpdatePorts      string
	appUpdateDomains    string
)

var appsUpdateCmd = &cobra.Command{
	Use:               "update <uuid|nome>",
	Short:             "Atualiza configurações de uma aplicação",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeApps,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()

		req := api.UpdateAppRequest{}
		if cmd.Flags().Changed("name") {
			v := appUpdateName
			req.Name = &v
		}
		if cmd.Flags().Changed("build-pack") {
			v := appUpdateBuildPack
			req.BuildPack = &v
		}
		if cmd.Flags().Changed("branch") {
			v := appUpdateBranch
			req.GitBranch = &v
		}
		if cmd.Flags().Changed("repository") {
			v := appUpdateRepository
			req.GitRepository = &v
		}
		if cmd.Flags().Changed("ports") {
			v := appUpdatePorts
			req.PortsExposes = &v
		}
		if cmd.Flags().Changed("domains") {
			v := appUpdateDomains
			req.FQDN = &v
		}

		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("Atualizando %q", args[0]))
		app, err := client.UpdateApp(args[0], req)
		if err != nil {
			spin.Fail("falhou")
			logger.Log(logger.ActionAppUpdate, logger.ResourceApp, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop("Aplicação atualizada")
		logger.Log(logger.ActionAppUpdate, logger.ResourceApp, args[0], app.Name, "success", time.Since(start))

		fmt.Println()
		kv := ui.NewTable("", "")
		kv.AddRow(ui.Bold("UUID:"), app.UUID)
		kv.AddRow(ui.Bold("Nome:"), app.Name)
		kv.AddRow(ui.Bold("Build pack:"), app.BuildPack)
		kv.AddRow(ui.Bold("Branch:"), app.Branch)
		kv.AddRow(ui.Bold("Repo:"), app.Repository)
		kv.AddRow(ui.Bold("Portas:"), app.PortsExposes)
		kv.AddRow(ui.Bold("Domínios:"), app.Domains)
		kv.Render()
		return nil
	},
}

func init() {
	appsUpdateCmd.Flags().StringVar(&appUpdateName, "name", "", "Novo nome")
	appsUpdateCmd.Flags().StringVar(&appUpdateBuildPack, "build-pack", "", "Build pack (nixpacks, dockerfile, static, dockercompose)")
	appsUpdateCmd.Flags().StringVar(&appUpdateBranch, "branch", "", "Branch do git")
	appsUpdateCmd.Flags().StringVar(&appUpdateRepository, "repository", "", "Repositório (owner/repo)")
	appsUpdateCmd.Flags().StringVar(&appUpdatePorts, "ports", "", "Portas expostas")
	appsUpdateCmd.Flags().StringVar(&appUpdateDomains, "domains", "", "Domínios (separados por vírgula)")
	appsCmd.AddCommand(appsUpdateCmd)
}
