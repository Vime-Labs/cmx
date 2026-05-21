package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/ui"
)

var dbTypes = []string{
	"postgresql", "mysql", "mariadb", "mongodb",
	"redis", "dragonfly", "keydb", "clickhouse",
}

var dbsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Cria um novo banco de dados (wizard interativo)",
	RunE:  runDBsCreate,
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
	fmt.Println()

	if len(projects) == 0 {
		return fmt.Errorf("nenhum projeto encontrado — crie um projeto no painel Coolify primeiro")
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
	name, err := ui.Input("Nome do banco", "", func(v string) error {
		if strings.ContainsAny(v, " \t/\\") {
			return fmt.Errorf("nome não pode conter espaços ou barras")
		}
		return nil
	})
	if err != nil {
		return err
	}

	// ── 6. Imagem Docker ──────────────────────────────────────────────────────
	defaultImage := api.DBDefaultImages[dbType]
	image, err := ui.Input("Imagem Docker", defaultImage, func(v string) error {
		if !strings.Contains(v, ":") {
			return fmt.Errorf("inclua a tag da versão (ex: %s)", defaultImage)
		}
		return nil
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

	if !ui.Confirm("Criar banco de dados?") {
		fmt.Println("Cancelado.")
		return nil
	}

	// ── 9. Criar ──────────────────────────────────────────────────────────────
	fmt.Print("\nCriando banco de dados... ")
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
		return err
	}
	fmt.Println("OK")
	fmt.Printf("\nBanco criado: %s\n", resp.UUID)
	fmt.Printf("\nVerifique o status com: cmx dbs get %s\n", name)
	return nil
}

func init() {
	dbsCmd.AddCommand(dbsCreateCmd)
}
