package cmd

import (
	"fmt"
	"time"

	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var (
	deployTag         string
	deployForceGlobal bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy [uuid|nome]",
	Short: "Dispara deploy por UUID, nome ou tag",
	Long: `Dispara o deploy de uma ou mais aplicações.

Exemplos:
  cmx deploy meu-app
  cmx deploy a1b2c3d4
  cmx deploy --tag production
  cmx deploy meu-app --force`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && deployTag == "" {
			return fmt.Errorf("informe um uuid/nome ou use --tag <tag>")
		}

		client := mustClient()

		if deployTag != "" {
			start := time.Now()
			spin := ui.NewSpinner(fmt.Sprintf("Disparando deploy para tag %q", deployTag))
			resp, err := client.DeployByTag(deployTag, deployForceGlobal)
			if err != nil {
				spin.Fail(err.Error())
				logger.Log(logger.ActionDeploy, logger.ResourceApp, deployTag, err.Error(), "error", time.Since(start))
				return err
			}
			spin.Stop(fmt.Sprintf("%d deploy(s) enfileirado(s)", len(resp.Deployments)))
			logger.Log(logger.ActionDeploy, logger.ResourceApp, deployTag,
				fmt.Sprintf("%d deployment(s)", len(resp.Deployments)), "success", time.Since(start))
			for _, d := range resp.Deployments {
				ui.Info(fmt.Sprintf("%s → %s", d.ResourceUUID, d.DeploymentUUID))
			}
			return nil
		}

		start := time.Now()
		spin := ui.NewSpinner(fmt.Sprintf("Disparando deploy para %q", args[0]))
		resp, err := client.DeployApp(args[0], deployForceGlobal)
		if err != nil {
			spin.Fail(err.Error())
			logger.Log(logger.ActionDeploy, logger.ResourceApp, args[0], err.Error(), "error", time.Since(start))
			return err
		}
		spin.Stop("Deploy enfileirado")
		logger.Log(logger.ActionDeploy, logger.ResourceApp, args[0],
			fmt.Sprintf("%d deployment(s)", len(resp.Deployments)), "success", time.Since(start))
		for _, d := range resp.Deployments {
			ui.Info(fmt.Sprintf("deployment: %s", d.DeploymentUUID))
		}
		return nil
	},
}

func init() {
	deployCmd.Flags().StringVar(&deployTag, "tag", "", "Dispara deploy em todos os recursos com esta tag")
	deployCmd.Flags().BoolVarP(&deployForceGlobal, "force", "f", false, "Força rebuild sem cache")
	rootCmd.AddCommand(deployCmd)
}
