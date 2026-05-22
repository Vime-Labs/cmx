package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/Vime-Labs/cmx/internal/validate"
	"github.com/spf13/cobra"
)

var validBuildPacks = []string{"nixpacks", "dockerfile", "static", "dockercompose"}

// ── flags ──────────────────────────────────────────────────────────────────────
var (
	appsCreateProject     string
	appsCreateEnvironment string
	appsCreateServer      string
	appsCreateGitHubApp   string
	appsCreateRepository  string
	appsCreateBranch      string
	appsCreateBuildPack   string
	appsCreatePorts       string
	appsCreateName        string
	appsCreateDomains     string
	appsCreateYes         bool
)

var appsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Cria uma nova aplicação (wizard interativo)",
	RunE:  runAppsCreate,
}

// ── helpers de resolução ───────────────────────────────────────────────────────

// isAppsNonInteractive retorna true se o usuário forneceu os flags essenciais
// para criar uma aplicação sem o wizard interativo.
func isAppsNonInteractive() bool {
	return appsCreateProject != "" &&
		appsCreateEnvironment != "" &&
		appsCreateServer != "" &&
		appsCreateGitHubApp != "" &&
		appsCreateRepository != ""
}

func resolveProject(projects []api.Project, id string) (*api.Project, error) {
	for _, p := range projects {
		if strings.EqualFold(p.UUID, id) || strings.EqualFold(p.Name, id) {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("projeto %q não encontrado", id)
}

func resolveEnvironment(envs []api.Environment, id string) (*api.Environment, error) {
	for _, e := range envs {
		if strings.EqualFold(e.UUID, id) || strings.EqualFold(e.Name, id) {
			return &e, nil
		}
	}
	return nil, fmt.Errorf("ambiente %q não encontrado", id)
}

func resolveServer(servers []api.Server, id string) (*api.Server, error) {
	for _, s := range servers {
		if strings.EqualFold(s.UUID, id) || strings.EqualFold(s.Name, id) {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("servidor %q não encontrado", id)
}

func resolveGitHubApp(ghApps []api.GitHubApp, id string) (*api.GitHubApp, error) {
	for _, g := range ghApps {
		if strings.EqualFold(g.UUID, id) || strings.EqualFold(g.Name, id) {
			return &g, nil
		}
	}
	return nil, fmt.Errorf("GitHub App %q não encontrado", id)
}

// ── runAppsCreate ──────────────────────────────────────────────────────────────

func runAppsCreate(cmd *cobra.Command, args []string) error {
	client := mustClient()

	// ── pré-carrega recursos ──────────────────────────────────────────────────
	spin := ui.NewSpinner("Buscando recursos do Coolify")
	projects, err := client.ListProjects()
	if err != nil {
		spin.Fail("falhou")
		return fmt.Errorf("buscando projetos: %w", err)
	}
	servers, err := client.ListServers()
	if err != nil {
		spin.Fail("falhou")
		return fmt.Errorf("buscando servidores: %w", err)
	}
	ghApps, err := client.ListGitHubApps()
	if err != nil {
		spin.Fail("falhou")
		return fmt.Errorf("buscando GitHub Apps: %w", err)
	}
	spin.Stop("Pronto")

	// ── Modo não-interativo ───────────────────────────────────────────────────
	if isAppsNonInteractive() {
		return runAppsCreateNonInteractive(client, projects, servers, ghApps)
	}

	// ── Modo interativo (wizard) ──────────────────────────────────────────────
	return runAppsCreateInteractive(client, projects, servers, ghApps)
}

// ── runAppsCreateNonInteractive ────────────────────────────────────────────────

func runAppsCreateNonInteractive(client api.API, projects []api.Project, servers []api.Server, ghApps []api.GitHubApp) error {
	// Resolve project
	proj, err := resolveProject(projects, appsCreateProject)
	if err != nil {
		return err
	}

	// Resolve environment within project
	fullProj, err := client.GetProject(proj.UUID)
	if err != nil {
		return fmt.Errorf("buscando ambientes do projeto: %w", err)
	}
	env, err := resolveEnvironment(fullProj.Environments, appsCreateEnvironment)
	if err != nil {
		return err
	}

	// Resolve server
	srv, err := resolveServer(servers, appsCreateServer)
	if err != nil {
		return err
	}

	// Resolve GitHub App
	gh, err := resolveGitHubApp(ghApps, appsCreateGitHubApp)
	if err != nil {
		return err
	}

	repo := appsCreateRepository
	branch := appsCreateBranch
	if branch == "" {
		branch = "main"
	}
	buildPack := appsCreateBuildPack
	if buildPack == "" {
		buildPack = "nixpacks"
	}
	ports := appsCreatePorts
	if ports == "" {
		ports = "3000"
	}
	name := appsCreateName
	if name == "" {
		name = repo[strings.Index(repo, "/")+1:]
	}
	domains := appsCreateDomains

	start := time.Now()
	fmt.Print("\nCriando aplicação... ")
	resp, err := client.CreateApp(api.CreateAppRequest{
		ProjectUUID:     proj.UUID,
		ServerUUID:      srv.UUID,
		EnvironmentName: env.Name,
		GitHubAppUUID:   gh.UUID,
		GitRepository:   repo,
		GitBranch:       branch,
		BuildPack:       buildPack,
		PortsExposes:    ports,
		Name:            name,
		FQDN:            domains,
	})
	if err != nil {
		fmt.Println("falhou")
		logger.Log(logger.ActionAppCreate, logger.ResourceApp, name, err.Error(), "error", time.Since(start))
		return err
	}
	fmt.Println("OK")
	logger.Log(logger.ActionAppCreate, logger.ResourceApp, name,
		fmt.Sprintf("UUID: %s", resp.UUID), "success", time.Since(start))
	fmt.Printf("\nAplicação criada: %s\n", resp.UUID)
	if resp.DeploymentUUID != "" {
		fmt.Printf("Deploy iniciado:  %s\n", resp.DeploymentUUID)
	}
	return nil
}

// ── runAppsCreateInteractive ───────────────────────────────────────────────────

func runAppsCreateInteractive(client api.API, projects []api.Project, servers []api.Server, ghApps []api.GitHubApp) error {
	fmt.Println()

	if len(projects) == 0 {
		fmt.Println("Nenhum projeto encontrado.")
		if !ui.Confirm("Criar um novo projeto?") {
			return fmt.Errorf("é necessário um projeto para criar a aplicação")
		}
		name, err := ui.Input("Nome do projeto", "meu-projeto", nil)
		if err != nil {
			return err
		}
		proj, err := client.CreateProject(name)
		if err != nil {
			return fmt.Errorf("criando projeto: %w", err)
		}
		projects = []api.Project{*proj}
	}
	if len(servers) == 0 {
		return fmt.Errorf("nenhum servidor encontrado — adicione um servidor no painel Coolify primeiro")
	}
	if len(ghApps) == 0 {
		return fmt.Errorf("nenhum GitHub App configurado — adicione um em: Configurações → GitHub Apps")
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
	repo, err := ui.Input("Repositório (owner/repo)", "", validate.RepoFormat)
	if err != nil {
		return err
	}

	// ── 6. Branch ─────────────────────────────────────────────────────────────
	branch, err := ui.Input("Branch", "main", validate.BranchName)
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
	ports, err := ui.Input("Porta(s) expostas", "3000", validate.Ports)
	if err != nil {
		return err
	}

	// ── 9. Nome ───────────────────────────────────────────────────────────────
	defaultName := repo[strings.Index(repo, "/")+1:]
	name, err := ui.Input("Nome da aplicação", defaultName, validate.ResourceName)
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

	if !appsCreateYes && !ui.Confirm("Criar aplicação?") {
		fmt.Println("Cancelado.")
		return nil
	}

	// ── 12. Criar ─────────────────────────────────────────────────────────────
	fmt.Print("\nCriando aplicação... ")
	start := time.Now()
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
		logger.Log(logger.ActionAppCreate, logger.ResourceApp, name, err.Error(), "error", time.Since(start))
		return err
	}
	fmt.Println("OK")
	logger.Log(logger.ActionAppCreate, logger.ResourceApp, name,
		fmt.Sprintf("UUID: %s", resp.UUID), "success", time.Since(start))
	fmt.Printf("\nAplicação criada: %s\n", resp.UUID)
	if resp.DeploymentUUID != "" {
		fmt.Printf("Deploy iniciado:  %s\n", resp.DeploymentUUID)
		fmt.Printf("\nAcompanhe com: cmx apps logs %s\n", name)
	}
	return nil
}

func init() {
	appsCreateCmd.Flags().StringVarP(&appsCreateProject, "project", "p", "", "Nome ou UUID do projeto")
	appsCreateCmd.Flags().StringVarP(&appsCreateEnvironment, "environment", "e", "", "Nome do ambiente")
	appsCreateCmd.Flags().StringVarP(&appsCreateServer, "server", "s", "", "Nome ou UUID do servidor")
	appsCreateCmd.Flags().StringVarP(&appsCreateGitHubApp, "github-app", "g", "", "Nome ou UUID do GitHub App")
	appsCreateCmd.Flags().StringVarP(&appsCreateRepository, "repository", "r", "", "Repositório (owner/repo)")
	appsCreateCmd.Flags().StringVarP(&appsCreateBranch, "branch", "b", "", "Branch (default: main)")
	appsCreateCmd.Flags().StringVar(&appsCreateBuildPack, "build-pack", "", "Build pack (nixpacks, dockerfile, static, dockercompose)")
	appsCreateCmd.Flags().StringVar(&appsCreatePorts, "ports", "", "Portas expostas (default: 3000)")
	appsCreateCmd.Flags().StringVarP(&appsCreateName, "name", "n", "", "Nome da aplicação (default: nome do repo)")
	appsCreateCmd.Flags().StringVar(&appsCreateDomains, "domains", "", "Domínios")
	appsCreateCmd.Flags().BoolVarP(&appsCreateYes, "yes", "y", false, "Pula confirmação")
	appsCmd.AddCommand(appsCreateCmd)
}
