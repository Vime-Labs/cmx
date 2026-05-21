package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/ui"
)

var validBuildPacks = []string{"nixpacks", "dockerfile", "static", "dockercompose"}

var appsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Cria uma nova aplicação (wizard interativo)",
	RunE:  runAppsCreate,
}

func runAppsCreate(cmd *cobra.Command, args []string) error {
	client := mustClient()

	// ── 1. Projeto ────────────────────────────────────────────────────────────
	projects, err := client.ListProjects()
	if err != nil {
		return fmt.Errorf("buscando projetos: %w", err)
	}
	if len(projects) == 0 {
		return fmt.Errorf("nenhum projeto encontrado — crie um projeto no painel Coolify primeiro")
	}
	projNames := make([]string, len(projects))
	for i, p := range projects {
		projNames[i] = p.Name
	}
	projIdx, err := ui.Select("Projeto", projNames)
	if err != nil {
		return err
	}
	project := projects[projIdx]

	// ── 2. Ambiente ───────────────────────────────────────────────────────────
	proj, err := client.GetProject(project.UUID)
	if err != nil {
		return fmt.Errorf("buscando ambientes do projeto: %w", err)
	}
	if len(proj.Environments) == 0 {
		return fmt.Errorf("projeto %q não tem ambientes configurados", project.Name)
	}
	envNames := make([]string, len(proj.Environments))
	for i, e := range proj.Environments {
		envNames[i] = e.Name
	}
	envIdx, err := ui.Select("Ambiente", envNames)
	if err != nil {
		return err
	}
	environment := proj.Environments[envIdx]

	// ── 3. Servidor ───────────────────────────────────────────────────────────
	servers, err := client.ListServers()
	if err != nil {
		return fmt.Errorf("buscando servidores: %w", err)
	}
	if len(servers) == 0 {
		return fmt.Errorf("nenhum servidor encontrado — adicione um servidor no painel Coolify primeiro")
	}
	srvNames := make([]string, len(servers))
	for i, s := range servers {
		srvNames[i] = fmt.Sprintf("%s (%s)", s.Name, s.IP)
	}
	srvIdx, err := ui.Select("Servidor", srvNames)
	if err != nil {
		return err
	}
	server := servers[srvIdx]

	// ── 4. GitHub App ─────────────────────────────────────────────────────────
	ghApps, err := client.ListGitHubApps()
	if err != nil {
		return fmt.Errorf("buscando GitHub Apps: %w", err)
	}
	if len(ghApps) == 0 {
		return fmt.Errorf("nenhum GitHub App configurado — adicione um em: Configurações → GitHub Apps")
	}
	ghNames := make([]string, len(ghApps))
	for i, g := range ghApps {
		ghNames[i] = g.Name
	}
	ghIdx, err := ui.Select("GitHub App", ghNames)
	if err != nil {
		return err
	}
	ghApp := ghApps[ghIdx]

	// ── 5. Repositório ────────────────────────────────────────────────────────
	repoRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)
	repo, err := ui.Input("Repositório (owner/repo)", "", func(v string) error {
		if !repoRegex.MatchString(v) {
			return fmt.Errorf("formato inválido — use owner/repo (ex: Vime-Labs/meu-app)")
		}
		return nil
	})
	if err != nil {
		return err
	}

	// ── 6. Branch ─────────────────────────────────────────────────────────────
	branch, err := ui.Input("Branch", "main", func(v string) error {
		if strings.ContainsAny(v, " \t") {
			return fmt.Errorf("branch não pode conter espaços")
		}
		return nil
	})
	if err != nil {
		return err
	}

	// ── 7. Build pack ─────────────────────────────────────────────────────────
	bpIdx, err := ui.Select("Build pack", validBuildPacks)
	if err != nil {
		return err
	}
	buildPack := validBuildPacks[bpIdx]

	// ── 8. Porta(s) ───────────────────────────────────────────────────────────
	ports, err := ui.Input("Porta(s) expostas", "3000", validatePorts)
	if err != nil {
		return err
	}

	// ── 9. Nome ───────────────────────────────────────────────────────────────
	defaultName := repo[strings.Index(repo, "/")+1:]
	name, err := ui.Input("Nome da aplicação", defaultName, nil)
	if err != nil {
		return err
	}

	// ── 10. Domínios ──────────────────────────────────────────────────────────
	domains := ui.InputOptional("Domínio(s)")

	// ── 11. Resumo + confirmação ──────────────────────────────────────────────
	ui.Summary([][2]string{
		{"Projeto:", project.Name},
		{"Ambiente:", environment.Name},
		{"Servidor:", servers[srvIdx].Name},
		{"GitHub App:", ghApp.Name},
		{"Repositório:", repo},
		{"Branch:", branch},
		{"Build pack:", buildPack},
		{"Porta(s):", ports},
		{"Nome:", name},
		{"Domínios:", domains},
	})

	if !ui.Confirm("Criar aplicação?") {
		fmt.Println("Cancelado.")
		return nil
	}

	// ── 12. Criar ─────────────────────────────────────────────────────────────
	fmt.Print("\nCriando aplicação... ")
	resp, err := client.CreateApp(api.CreateAppRequest{
		ProjectUUID:     project.UUID,
		ServerUUID:      server.UUID,
		EnvironmentName: environment.Name,
		GitHubAppUUID:   ghApp.UUID,
		GitRepository:   repo,
		GitBranch:       branch,
		BuildPack:       buildPack,
		PortsExposes:    ports,
		Name:            name,
		FQDN:            domains,
	})
	if err != nil {
		fmt.Println("falhou")
		return err
	}
	fmt.Println("OK")
	fmt.Printf("\nAplicação criada: %s\n", resp.UUID)
	if resp.DeploymentUUID != "" {
		fmt.Printf("Deploy iniciado:  %s\n", resp.DeploymentUUID)
		fmt.Printf("\nAcompanhe com: cmx apps logs %s\n", name)
	}
	return nil
}

// validatePorts aceita "3000" ou "3000,3001" — cada parte deve ser 1-65535.
func validatePorts(v string) error {
	parts := strings.Split(v, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 65535 {
			return fmt.Errorf("%q não é uma porta válida (1-65535)", p)
		}
	}
	return nil
}

func init() {
	appsCmd.AddCommand(appsCreateCmd)
}
