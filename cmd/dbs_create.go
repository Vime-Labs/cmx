package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/logger"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/Vime-Labs/cmx/internal/validate"
	"github.com/spf13/cobra"
)

var dbTypes = []string{
	"postgresql", "mysql", "mariadb", "mongodb",
	"redis", "dragonfly", "keydb", "clickhouse",
}

// ── flags ──────────────────────────────────────────────────────────────────────
var (
	dbsCreateProject     string
	dbsCreateEnvironment string
	dbsCreateServer      string
	dbsCreateType        string
	dbsCreateName        string
	dbsCreateImage       string
	dbsCreatePublic      bool
	dbsCreatePublicPort  int
	dbsCreateYes         bool
)

var dbsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Cria um novo banco de dados (wizard interativo)",
	RunE:  runDBsCreate,
}

// isDbsNonInteractive retorna true se o usuário forneceu os flags essenciais
// para criar um banco sem o wizard interativo.
func isDbsNonInteractive() bool {
	return dbsCreateProject != "" &&
		dbsCreateEnvironment != "" &&
		dbsCreateServer != "" &&
		dbsCreateType != "" &&
		dbsCreateName != ""
}

// anyDBFlagSet retorna true se algum flag de criação foi fornecido.
func anyDBFlagSet() bool {
	return dbsCreateProject != "" || dbsCreateEnvironment != "" ||
		dbsCreateServer != "" || dbsCreateType != "" ||
		dbsCreateName != "" || dbsCreateImage != "" ||
		dbsCreatePublic || dbsCreatePublicPort != 0
}

// missingRequiredDBFlags retorna os nomes dos flags obrigatórios que faltam.
func missingRequiredDBFlags() []string {
	var missing []string
	if dbsCreateProject == "" {
		missing = append(missing, "--project")
	}
	if dbsCreateEnvironment == "" {
		missing = append(missing, "--environment")
	}
	if dbsCreateServer == "" {
		missing = append(missing, "--server")
	}
	if dbsCreateType == "" {
		missing = append(missing, "--type")
	}
	if dbsCreateName == "" {
		missing = append(missing, "--name")
	}
	return missing
}

func runDBsCreate(cmd *cobra.Command, args []string) error {
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
	spin.Stop("Pronto")

	// ── Modo não-interativo ───────────────────────────────────────────────────
	if isDbsNonInteractive() {
		return runDBsCreateNonInteractive(client, projects, servers)
	}

	// ── Modo parcial: alguns flags fornecidos mas não todos ───────────────────
	if anyDBFlagSet() {
		missing := missingRequiredDBFlags()
		return fmt.Errorf("modo nao-interativo requer todos os flags obrigatorios. Faltam: %s",
			strings.Join(missing, ", "))
	}

	// ── Modo interativo (wizard) ──────────────────────────────────────────────
	return runDBsCreateInteractive(client, projects, servers)
}

// runDBsCreateNonInteractive cria o banco com os valores fornecidos via flags.
func runDBsCreateNonInteractive(client api.API, projects []api.Project, servers []api.Server) error {
	// Resolve project
	proj, err := resolveProject(projects, dbsCreateProject)
	if err != nil {
		return err
	}

	// Resolve environment within project
	fullProj, err := client.GetProject(proj.UUID)
	if err != nil {
		return fmt.Errorf("buscando ambientes do projeto: %w", err)
	}
	env, err := resolveEnvironment(fullProj.Environments, dbsCreateEnvironment)
	if err != nil {
		return err
	}

	// Resolve server
	srv, err := resolveServer(servers, dbsCreateServer)
	if err != nil {
		return err
	}

	dbType := dbsCreateType
	name := dbsCreateName
	image := dbsCreateImage
	if image == "" {
		if def, ok := api.DBDefaultImages[dbType]; ok {
			image = def
		}
	}
	isPublic := dbsCreatePublic
	publicPort := dbsCreatePublicPort

	start := time.Now()
	fmt.Print("\nCriando banco de dados... ")
	resp, err := client.CreateDB(dbType, api.CreateDBRequest{
		ProjectUUID:     proj.UUID,
		ServerUUID:      srv.UUID,
		EnvironmentName: env.Name,
		Name:            name,
		Image:           image,
		IsPublic:        isPublic,
		PublicPort:      publicPort,
	})
	if err != nil {
		fmt.Println("falhou")
		logger.Log(logger.ActionDBCreate, logger.ResourceDB, name, err.Error(), "error", time.Since(start))
		return err
	}
	fmt.Println("OK")
	logger.Log(logger.ActionDBCreate, logger.ResourceDB, name,
		fmt.Sprintf("UUID: %s", resp.UUID), "success", time.Since(start))
	fmt.Printf("\nBanco criado: %s\n", resp.UUID)
	return nil
}

// runDBsCreateInteractive executa o wizard interativo original.
func runDBsCreateInteractive(client api.API, projects []api.Project, servers []api.Server) error {
	fmt.Println()

	if len(projects) == 0 {
		fmt.Println("Nenhum projeto encontrado.")
		if !ui.Confirm("Criar um novo projeto?") {
			return fmt.Errorf("é necessário um projeto para criar o banco de dados")
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

	// ── 4. Tipo de banco ──────────────────────────────────────────────────────
	dbTypeIdx, err := ui.Select("Tipo de banco", dbTypes)
	if err != nil {
		return err
	}
	dbType := dbTypes[dbTypeIdx]

	// ── 5. Nome ───────────────────────────────────────────────────────────────
	name, err := ui.Input("Nome do banco", "", validate.ResourceName)
	if err != nil {
		return err
	}

	// ── 6. Imagem Docker ──────────────────────────────────────────────────────
	defaultImage := api.DBDefaultImages[dbType]
	image, err := ui.Input("Imagem Docker", defaultImage, func(v string) error {
		return validate.ImageTag(v, defaultImage)
	})
	if err != nil {
		return err
	}

	// ── 7. Porta pública ──────────────────────────────────────────────────────
	var isPublic bool
	var publicPort int
	if ui.Confirm("Expor porta publicamente?") {
		isPublic = true
		portStr, err := ui.Input("Porta pública", "", func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil || n < 1 || n > 65535 {
				return fmt.Errorf("porta inválida (1-65535)")
			}
			return nil
		})
		if err != nil {
			return err
		}
		publicPort, _ = strconv.Atoi(portStr)
	}

	// ── 8. Resumo + confirmação ───────────────────────────────────────────────
	publicInfo := "não"
	if isPublic {
		publicInfo = fmt.Sprintf("sim (porta %d)", publicPort)
	}
	ui.Summary([][2]string{
		{"Projeto:", project.Name},
		{"Ambiente:", environment.Name},
		{"Servidor:", servers[srvIdx].Name},
		{"Tipo:", dbType},
		{"Nome:", name},
		{"Imagem:", image},
		{"Público:", publicInfo},
	})

	if !dbsCreateYes && !ui.Confirm("Criar banco de dados?") {
		fmt.Println("Cancelado.")
		return nil
	}

	// ── 9. Criar ──────────────────────────────────────────────────────────────
	fmt.Print("\nCriando banco de dados... ")
	start := time.Now()
	resp, err := client.CreateDB(dbType, api.CreateDBRequest{
		ProjectUUID:     project.UUID,
		ServerUUID:      server.UUID,
		EnvironmentName: environment.Name,
		Name:            name,
		Image:           image,
		IsPublic:        isPublic,
		PublicPort:      publicPort,
	})
	if err != nil {
		fmt.Println("falhou")
		logger.Log(logger.ActionDBCreate, logger.ResourceDB, name, err.Error(), "error", time.Since(start))
		return err
	}
	fmt.Println("OK")
	logger.Log(logger.ActionDBCreate, logger.ResourceDB, name,
		fmt.Sprintf("UUID: %s", resp.UUID), "success", time.Since(start))
	fmt.Printf("\nBanco criado: %s\n", resp.UUID)
	fmt.Printf("\nVerifique o status com: cmx dbs get %s\n", name)
	return nil
}

func init() {
	dbsCreateCmd.Flags().StringVarP(&dbsCreateProject, "project", "p", "", "Nome ou UUID do projeto")
	dbsCreateCmd.Flags().StringVarP(&dbsCreateEnvironment, "environment", "e", "", "Nome do ambiente")
	dbsCreateCmd.Flags().StringVarP(&dbsCreateServer, "server", "s", "", "Nome ou UUID do servidor")
	dbsCreateCmd.Flags().StringVarP(&dbsCreateType, "type", "t", "", "Tipo de banco (postgresql, mysql, redis, etc.)")
	dbsCreateCmd.Flags().StringVarP(&dbsCreateName, "name", "n", "", "Nome do banco")
	dbsCreateCmd.Flags().StringVar(&dbsCreateImage, "image", "", "Imagem Docker (default conforme o tipo)")
	dbsCreateCmd.Flags().BoolVar(&dbsCreatePublic, "public", false, "Expor porta publicamente")
	dbsCreateCmd.Flags().IntVar(&dbsCreatePublicPort, "public-port", 0, "Porta pública")
	dbsCreateCmd.Flags().BoolVarP(&dbsCreateYes, "yes", "y", false, "Pula confirmação")

	// Flag completion
	dbsCreateCmd.RegisterFlagCompletionFunc("project", completeProjects)
	dbsCreateCmd.RegisterFlagCompletionFunc("server", completeServers)

	dbsCmd.AddCommand(dbsCreateCmd)
}
