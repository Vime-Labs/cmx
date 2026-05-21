package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Vime-Labs/cmx/internal/api"
)

// jsonSchema monta um json.RawMessage de propriedades para o inputSchema.
func jsonSchema(props map[string]map[string]interface{}) json.RawMessage {
	data, _ := json.Marshal(props)
	return data
}

func (s *Server) registerTools() {
	s.add(toolDefinition{
		Name:        "cmx_ping",
		Description: "Testa a conexão com o servidor Coolify",
		InputSchema: inputSchema{Type: "object", Properties: jsonSchema(nil)},
	}, handlePing)

	s.add(toolDefinition{
		Name:        "cmx_list_apps",
		Description: "Lista todas as aplicações",
		InputSchema: inputSchema{Type: "object", Properties: jsonSchema(nil)},
	}, handleListApps)

	s.add(toolDefinition{
		Name:        "cmx_get_app",
		Description: "Exibe detalhes de uma aplicação",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"id": {"type": "string", "description": "UUID ou nome da aplicação"}}),
			Required:   []string{"id"},
		},
	}, handleGetApp)

	s.add(toolDefinition{
		Name:        "cmx_app_logs",
		Description: "Exibe logs de uma aplicação",
		InputSchema: inputSchema{
			Type: "object",
			Properties: jsonSchema(map[string]map[string]interface{}{
				"id":    {"type": "string", "description": "UUID ou nome da aplicação"},
				"lines": {"type": "integer", "description": "Número de linhas (default 100)", "default": 100},
			}),
			Required: []string{"id"},
		},
	}, handleAppLogs)

	s.add(toolDefinition{
		Name:        "cmx_deploy_app",
		Description: "Dispara um deploy para uma aplicação",
		InputSchema: inputSchema{
			Type: "object",
			Properties: jsonSchema(map[string]map[string]interface{}{
				"id":    {"type": "string", "description": "UUID ou nome da aplicação"},
				"force": {"type": "boolean", "description": "Força rebuild sem cache", "default": false},
			}),
			Required: []string{"id"},
		},
	}, handleDeployApp)

	s.add(toolDefinition{
		Name:        "cmx_start_app",
		Description: "Inicia uma aplicação",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"id": {"type": "string", "description": "UUID ou nome da aplicação"}}),
			Required:   []string{"id"},
		},
	}, handleStartApp)

	s.add(toolDefinition{
		Name:        "cmx_stop_app",
		Description: "Para uma aplicação",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"id": {"type": "string", "description": "UUID ou nome da aplicação"}}),
			Required:   []string{"id"},
		},
	}, handleStopApp)

	s.add(toolDefinition{
		Name:        "cmx_restart_app",
		Description: "Reinicia uma aplicação",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"id": {"type": "string", "description": "UUID ou nome da aplicação"}}),
			Required:   []string{"id"},
		},
	}, handleRestartApp)

	s.add(toolDefinition{
		Name:        "cmx_create_app",
		Description: "Cria uma nova aplicação",
		InputSchema: inputSchema{
			Type: "object",
			Properties: jsonSchema(map[string]map[string]interface{}{
				"project_uuid":     {"type": "string", "description": "UUID do projeto"},
				"server_uuid":      {"type": "string", "description": "UUID do servidor"},
				"environment_name": {"type": "string", "description": "Nome do ambiente (ex: production, staging)"},
				"github_app_uuid":  {"type": "string", "description": "UUID do GitHub App"},
				"git_repository":   {"type": "string", "description": "Repositório no formato owner/repo"},
				"git_branch":       {"type": "string", "description": "Branch do git"},
				"build_pack":       {"type": "string", "description": "Build pack: nixpacks, dockerfile, static, dockercompose"},
				"ports_exposes":    {"type": "string", "description": "Portas expostas (ex: 3000 ou 3000,8080)"},
				"name":             {"type": "string", "description": "Nome da aplicação (opcional, default = nome do repo)"},
				"fqdn":             {"type": "string", "description": "Domínio(s) (opcional)"},
			}),
			Required: []string{"project_uuid", "server_uuid", "environment_name", "github_app_uuid", "git_repository", "git_branch", "build_pack", "ports_exposes"},
		},
	}, handleCreateApp)

	s.add(toolDefinition{
		Name:        "cmx_list_app_envs",
		Description: "Lista variáveis de ambiente de uma aplicação",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"id": {"type": "string", "description": "UUID ou nome da aplicação"}}),
			Required:   []string{"id"},
		},
	}, handleListAppEnvs)

	s.add(toolDefinition{
		Name:        "cmx_set_app_env",
		Description: "Cria ou atualiza uma variável de ambiente",
		InputSchema: inputSchema{
			Type: "object",
			Properties: jsonSchema(map[string]map[string]interface{}{
				"id":    {"type": "string", "description": "UUID ou nome da aplicação"},
				"key":   {"type": "string", "description": "Nome da variável"},
				"value": {"type": "string", "description": "Valor da variável"},
			}),
			Required: []string{"id", "key", "value"},
		},
	}, handleSetAppEnv)

	s.add(toolDefinition{
		Name:        "cmx_delete_app_env",
		Description: "Remove uma variável de ambiente de uma aplicação",
		InputSchema: inputSchema{
			Type: "object",
			Properties: jsonSchema(map[string]map[string]interface{}{
				"id":  {"type": "string", "description": "UUID ou nome da aplicação"},
				"key": {"type": "string", "description": "Nome da variável a remover"},
			}),
			Required: []string{"id", "key"},
		},
	}, handleDeleteAppEnv)

	s.add(toolDefinition{
		Name:        "cmx_list_dbs",
		Description: "Lista todos os bancos de dados",
		InputSchema: inputSchema{Type: "object", Properties: jsonSchema(nil)},
	}, handleListDBs)

	s.add(toolDefinition{
		Name:        "cmx_get_db",
		Description: "Exibe detalhes de um banco de dados",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"id": {"type": "string", "description": "UUID ou nome do banco"}}),
			Required:   []string{"id"},
		},
	}, handleGetDB)

	s.add(toolDefinition{
		Name:        "cmx_start_db",
		Description: "Inicia um banco de dados",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"id": {"type": "string", "description": "UUID ou nome do banco"}}),
			Required:   []string{"id"},
		},
	}, handleStartDB)

	s.add(toolDefinition{
		Name:        "cmx_stop_db",
		Description: "Para um banco de dados",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"id": {"type": "string", "description": "UUID ou nome do banco"}}),
			Required:   []string{"id"},
		},
	}, handleStopDB)

	s.add(toolDefinition{
		Name:        "cmx_restart_db",
		Description: "Reinicia um banco de dados",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"id": {"type": "string", "description": "UUID ou nome do banco"}}),
			Required:   []string{"id"},
		},
	}, handleRestartDB)

	s.add(toolDefinition{
		Name:        "cmx_create_db",
		Description: "Cria um novo banco de dados",
		InputSchema: inputSchema{
			Type: "object",
			Properties: jsonSchema(map[string]map[string]interface{}{
				"db_type":          {"type": "string", "description": "Tipo: postgresql, mysql, mariadb, mongodb, redis, dragonfly, keydb, clickhouse"},
				"project_uuid":     {"type": "string", "description": "UUID do projeto"},
				"server_uuid":      {"type": "string", "description": "UUID do servidor"},
				"environment_name": {"type": "string", "description": "Nome do ambiente"},
				"name":             {"type": "string", "description": "Nome do banco"},
				"image":            {"type": "string", "description": "Imagem Docker (opcional, ex: postgres:16)"},
				"is_public":        {"type": "boolean", "description": "Expor porta pública", "default": false},
				"public_port":      {"type": "integer", "description": "Porta pública (obrigatório se is_public=true)"},
			}),
			Required: []string{"db_type", "project_uuid", "server_uuid", "environment_name", "name"},
		},
	}, handleCreateDB)

	s.add(toolDefinition{
		Name:        "cmx_list_projects",
		Description: "Lista todos os projetos",
		InputSchema: inputSchema{Type: "object", Properties: jsonSchema(nil)},
	}, handleListProjects)

	s.add(toolDefinition{
		Name:        "cmx_list_servers",
		Description: "Lista todos os servidores",
		InputSchema: inputSchema{Type: "object", Properties: jsonSchema(nil)},
	}, handleListServers)

	s.add(toolDefinition{
		Name:        "cmx_list_github_apps",
		Description: "Lista GitHub Apps configurados",
		InputSchema: inputSchema{Type: "object", Properties: jsonSchema(nil)},
	}, handleListGitHubApps)

	s.add(toolDefinition{
		Name:        "cmx_list_deployments",
		Description: "Lista deployments em andamento",
		InputSchema: inputSchema{Type: "object", Properties: jsonSchema(nil)},
	}, handleListDeployments)

	s.add(toolDefinition{
		Name:        "cmx_get_deployment",
		Description: "Exibe detalhes de um deployment",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"uuid": {"type": "string", "description": "UUID do deployment"}}),
			Required:   []string{"uuid"},
		},
	}, handleGetDeployment)

	s.add(toolDefinition{
		Name:        "cmx_cancel_deployment",
		Description: "Cancela um deployment em andamento",
		InputSchema: inputSchema{
			Type:       "object",
			Properties: jsonSchema(map[string]map[string]interface{}{"uuid": {"type": "string", "description": "UUID do deployment"}}),
			Required:   []string{"uuid"},
		},
	}, handleCancelDeployment)

	s.add(toolDefinition{
		Name:        "cmx_list_app_deployments",
		Description: "Histórico de deployments de uma aplicação",
		InputSchema: inputSchema{
			Type: "object",
			Properties: jsonSchema(map[string]map[string]interface{}{
				"id":   {"type": "string", "description": "UUID ou nome da aplicação"},
				"skip": {"type": "integer", "description": "Pular N registros", "default": 0},
				"take": {"type": "integer", "description": "Exibir N registros", "default": 10},
			}),
			Required: []string{"id"},
		},
	}, handleListAppDeployments)

	s.add(toolDefinition{
		Name:        "cmx_deploy_by_tag",
		Description: "Dispara deploy em todos os recursos com uma tag",
		InputSchema: inputSchema{
			Type: "object",
			Properties: jsonSchema(map[string]map[string]interface{}{
				"tag":   {"type": "string", "description": "Tag para deploy"},
				"force": {"type": "boolean", "description": "Força rebuild sem cache", "default": false},
			}),
			Required: []string{"tag"},
		},
	}, handleDeployByTag)
}

// ── Handlers ─────────────────────────────────────────────────────────────────

func text(text string) callToolResult {
	return callToolResult{
		Content: []contentItem{{Type: "text", Text: text}},
	}
}

func errText(text string) callToolResult {
	return callToolResult{
		Content: []contentItem{{Type: "text", Text: text}},
		IsError: true,
	}
}

func handlePing(client api.API, _ json.RawMessage) callToolResult {
	if err := client.Ping(); err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text("Conexão OK")
}

func handleListApps(client api.API, _ json.RawMessage) callToolResult {
	apps, err := client.ListApps()
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	if len(apps) == 0 {
		return text("Nenhuma aplicação encontrada")
	}
	var b strings.Builder
	for _, a := range apps {
		fmt.Fprintf(&b, "%-10s  %-30s  %-12s  %s\n", a.UUID, a.Name, a.Status, a.Repository)
	}
	return text(b.String())
}

func handleGetApp(client api.API, params json.RawMessage) callToolResult {
	var p struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	app, err := client.GetApp(p.ID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "UUID:       %s\n", app.UUID)
	fmt.Fprintf(&b, "Nome:       %s\n", app.Name)
	fmt.Fprintf(&b, "Status:     %s\n", app.Status)
	fmt.Fprintf(&b, "Repo:       %s\n", app.Repository)
	fmt.Fprintf(&b, "Branch:     %s\n", app.Branch)
	fmt.Fprintf(&b, "BuildPack:  %s\n", app.BuildPack)
	fmt.Fprintf(&b, "Domínios:   %s\n", app.Domains)
	fmt.Fprintf(&b, "Criado em:  %s\n", app.CreatedAt)
	fmt.Fprintf(&b, "Atualizado: %s\n", app.UpdatedAt)
	return text(b.String())
}

func handleAppLogs(client api.API, params json.RawMessage) callToolResult {
	var p struct {
		ID    string `json:"id"`
		Lines int    `json:"lines"`
	}
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	if p.Lines <= 0 {
		p.Lines = 100
	}
	logs, err := client.AppLogs(p.ID, p.Lines)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(logs)
}

func handleDeployApp(client api.API, params json.RawMessage) callToolResult {
	var p struct {
		ID    string `json:"id"`
		Force bool   `json:"force"`
	}
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	resp, err := client.DeployApp(p.ID, p.Force)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Deploy enfileirado para %s:\n", p.ID)
	for _, d := range resp.Deployments {
		fmt.Fprintf(&b, "  deployment: %s\n", d.DeploymentUUID)
	}
	return text(b.String())
}

func handleStartApp(client api.API, params json.RawMessage) callToolResult {
	var p struct{ ID string `json:"id"` }
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	msg, err := client.StartApp(p.ID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(msg)
}

func handleStopApp(client api.API, params json.RawMessage) callToolResult {
	var p struct{ ID string `json:"id"` }
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	msg, err := client.StopApp(p.ID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(msg)
}

func handleRestartApp(client api.API, params json.RawMessage) callToolResult {
	var p struct{ ID string `json:"id"` }
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	msg, err := client.RestartApp(p.ID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(msg)
}

func handleCreateApp(client api.API, params json.RawMessage) callToolResult {
	var p struct {
		ProjectUUID     string `json:"project_uuid"`
		ServerUUID      string `json:"server_uuid"`
		EnvironmentName string `json:"environment_name"`
		GitHubAppUUID   string `json:"github_app_uuid"`
		GitRepository   string `json:"git_repository"`
		GitBranch       string `json:"git_branch"`
		BuildPack       string `json:"build_pack"`
		PortsExposes    string `json:"ports_exposes"`
		Name            string `json:"name"`
		FQDN            string `json:"fqdn"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return errText("Parâmetros inválidos")
	}
	req := api.CreateAppRequest{
		ProjectUUID:     p.ProjectUUID,
		ServerUUID:      p.ServerUUID,
		EnvironmentName: p.EnvironmentName,
		GitHubAppUUID:   p.GitHubAppUUID,
		GitRepository:   p.GitRepository,
		GitBranch:       p.GitBranch,
		BuildPack:       p.BuildPack,
		PortsExposes:    p.PortsExposes,
		Name:            p.Name,
		FQDN:            p.FQDN,
	}
	resp, err := client.CreateApp(req)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Aplicação criada: %s\n", resp.UUID)
	if resp.DeploymentUUID != "" {
		fmt.Fprintf(&b, "Deploy iniciado: %s\n", resp.DeploymentUUID)
	}
	return text(b.String())
}

func handleListAppEnvs(client api.API, params json.RawMessage) callToolResult {
	var p struct{ ID string `json:"id"` }
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	envs, err := client.ListAppEnvs(p.ID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	if len(envs) == 0 {
		return text("Nenhuma variável de ambiente encontrada")
	}
	var b strings.Builder
	for _, e := range envs {
		fmt.Fprintf(&b, "%-30s  %s", e.Key, maskSecret(e.Key, e.Value))
		if e.IsPreview {
			b.WriteString("  [preview]")
		}
		b.WriteString("\n")
	}
	return text(b.String())
}

func handleSetAppEnv(client api.API, params json.RawMessage) callToolResult {
	var p struct {
		ID    string `json:"id"`
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return errText("Parâmetros inválidos")
	}
	if p.ID == "" || p.Key == "" {
		return errText("Parâmetros \"id\" e \"key\" são obrigatórios")
	}
	if err := client.SetAppEnv(p.ID, p.Key, p.Value); err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(fmt.Sprintf("Variável %s definida", p.Key))
}

func handleDeleteAppEnv(client api.API, params json.RawMessage) callToolResult {
	var p struct {
		ID  string `json:"id"`
		Key string `json:"key"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return errText("Parâmetros inválidos")
	}
	if p.ID == "" || p.Key == "" {
		return errText("Parâmetros \"id\" e \"key\" são obrigatórios")
	}
	if err := client.DeleteAppEnvByKey(p.ID, p.Key); err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(fmt.Sprintf("Variável %s removida", p.Key))
}

func handleListDBs(client api.API, _ json.RawMessage) callToolResult {
	dbs, err := client.ListDBs()
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	if len(dbs) == 0 {
		return text("Nenhum banco de dados encontrado")
	}
	var b strings.Builder
	for _, d := range dbs {
		public := "-"
		if d.IsPublic {
			public = fmt.Sprintf("%d", d.PublicPort)
		}
		fmt.Fprintf(&b, "%-10s  %-25s  %-12s  %-20s  %-12s  %s\n", d.UUID, d.Name, d.DisplayType(), d.Image, d.Status, public)
	}
	return text(b.String())
}

func handleGetDB(client api.API, params json.RawMessage) callToolResult {
	var p struct{ ID string `json:"id"` }
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	db, err := client.GetDB(p.ID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "UUID:     %s\n", db.UUID)
	fmt.Fprintf(&b, "Nome:     %s\n", db.Name)
	fmt.Fprintf(&b, "Tipo:     %s\n", db.DisplayType())
	fmt.Fprintf(&b, "Imagem:   %s\n", db.Image)
	fmt.Fprintf(&b, "Status:   %s\n", db.Status)
	if db.IsPublic {
		fmt.Fprintf(&b, "Público:  sim (porta %d)\n", db.PublicPort)
	} else {
		fmt.Fprintf(&b, "Público:  não\n")
	}
	fmt.Fprintf(&b, "Criado:   %s\n", db.CreatedAt)
	fmt.Fprintf(&b, "Atualizado: %s\n", db.UpdatedAt)
	return text(b.String())
}

func handleStartDB(client api.API, params json.RawMessage) callToolResult {
	var p struct{ ID string `json:"id"` }
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	msg, err := client.StartDB(p.ID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(msg)
}

func handleStopDB(client api.API, params json.RawMessage) callToolResult {
	var p struct{ ID string `json:"id"` }
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	msg, err := client.StopDB(p.ID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(msg)
}

func handleRestartDB(client api.API, params json.RawMessage) callToolResult {
	var p struct{ ID string `json:"id"` }
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	msg, err := client.RestartDB(p.ID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(msg)
}

func handleCreateDB(client api.API, params json.RawMessage) callToolResult {
	var p struct {
		DBType          string `json:"db_type"`
		ProjectUUID     string `json:"project_uuid"`
		ServerUUID      string `json:"server_uuid"`
		EnvironmentName string `json:"environment_name"`
		Name            string `json:"name"`
		Image           string `json:"image"`
		IsPublic        bool   `json:"is_public"`
		PublicPort      int    `json:"public_port"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return errText("Parâmetros inválidos")
	}
	if p.DBType == "" || p.ProjectUUID == "" || p.ServerUUID == "" || p.EnvironmentName == "" || p.Name == "" {
		return errText("Parâmetros obrigatórios: db_type, project_uuid, server_uuid, environment_name, name")
	}
	if p.Image == "" {
		if img, ok := api.DBDefaultImages[p.DBType]; ok {
			p.Image = img
		}
	}
	req := api.CreateDBRequest{
		ProjectUUID:     p.ProjectUUID,
		ServerUUID:      p.ServerUUID,
		EnvironmentName: p.EnvironmentName,
		Name:            p.Name,
		Image:           p.Image,
		IsPublic:        p.IsPublic,
		PublicPort:      p.PublicPort,
	}
	resp, err := client.CreateDB(p.DBType, req)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text(fmt.Sprintf("Banco criado: %s\n%s", resp.UUID, resp.Message))
}

func handleListProjects(client api.API, _ json.RawMessage) callToolResult {
	projects, err := client.ListProjects()
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	if len(projects) == 0 {
		return text("Nenhum projeto encontrado")
	}
	var b strings.Builder
	for _, p := range projects {
		fmt.Fprintf(&b, "%-10s  %s", p.UUID, p.Name)
		if p.Description != "" {
			fmt.Fprintf(&b, "  (%s)", p.Description)
		}
		b.WriteString("\n")
	}
	return text(b.String())
}

func handleListServers(client api.API, _ json.RawMessage) callToolResult {
	servers, err := client.ListServers()
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	if len(servers) == 0 {
		return text("Nenhum servidor encontrado")
	}
	var b strings.Builder
	for _, s := range servers {
		fmt.Fprintf(&b, "%-10s  %-20s  %s\n", s.UUID, s.Name, s.IP)
	}
	return text(b.String())
}

func handleListGitHubApps(client api.API, _ json.RawMessage) callToolResult {
	apps, err := client.ListGitHubApps()
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	if len(apps) == 0 {
		return text("Nenhum GitHub App configurado")
	}
	var b strings.Builder
	for _, a := range apps {
		fmt.Fprintf(&b, "%-10s  %s\n", a.UUID, a.Name)
	}
	return text(b.String())
}

func handleListDeployments(client api.API, _ json.RawMessage) callToolResult {
	deps, err := client.ListDeployments()
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	if len(deps) == 0 {
		return text("Nenhum deployment ativo")
	}
	var b strings.Builder
	for _, d := range deps {
		fmt.Fprintf(&b, "%-10s  %-10s  %-12s  %s\n", d.UUID, d.ApplicationUUID, d.Status, d.CreatedAt)
	}
	return text(b.String())
}

func handleGetDeployment(client api.API, params json.RawMessage) callToolResult {
	var p struct{ UUID string `json:"uuid"` }
	if err := json.Unmarshal(params, &p); err != nil || p.UUID == "" {
		return errText("Parâmetro \"uuid\" é obrigatório")
	}
	d, err := client.GetDeployment(p.UUID)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "UUID:       %s\n", d.UUID)
	fmt.Fprintf(&b, "App UUID:   %s\n", d.ApplicationUUID)
	fmt.Fprintf(&b, "Status:     %s\n", d.Status)
	fmt.Fprintf(&b, "Criado em:  %s\n", d.CreatedAt)
	fmt.Fprintf(&b, "Atualizado: %s\n", d.UpdatedAt)
	return text(b.String())
}

func handleCancelDeployment(client api.API, params json.RawMessage) callToolResult {
	var p struct{ UUID string `json:"uuid"` }
	if err := json.Unmarshal(params, &p); err != nil || p.UUID == "" {
		return errText("Parâmetro \"uuid\" é obrigatório")
	}
	if err := client.CancelDeployment(p.UUID); err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	return text("Deployment cancelado")
}

func handleListAppDeployments(client api.API, params json.RawMessage) callToolResult {
	var p struct {
		ID   string `json:"id"`
		Skip int    `json:"skip"`
		Take int    `json:"take"`
	}
	if err := json.Unmarshal(params, &p); err != nil || p.ID == "" {
		return errText("Parâmetro \"id\" é obrigatório")
	}
	if p.Take <= 0 {
		p.Take = 10
	}
	deps, err := client.ListAppDeployments(p.ID, p.Skip, p.Take)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	if len(deps) == 0 {
		return text("Nenhum deployment encontrado")
	}
	var b strings.Builder
	for _, d := range deps {
		fmt.Fprintf(&b, "%-10s  %-12s  %s  %s\n", d.UUID, d.Status, d.CreatedAt, d.UpdatedAt)
	}
	return text(b.String())
}

func handleDeployByTag(client api.API, params json.RawMessage) callToolResult {
	var p struct {
		Tag   string `json:"tag"`
		Force bool   `json:"force"`
	}
	if err := json.Unmarshal(params, &p); err != nil || p.Tag == "" {
		return errText("Parâmetro \"tag\" é obrigatório")
	}
	resp, err := client.DeployByTag(p.Tag, p.Force)
	if err != nil {
		return errText(fmt.Sprintf("Erro: %v", err))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Deploy enfileirado para tag %q:\n", p.Tag)
	for _, d := range resp.Deployments {
		fmt.Fprintf(&b, "  %s → %s\n", d.ResourceUUID, d.DeploymentUUID)
	}
	return text(b.String())
}

// maskSecret oculta parcialmente valores de variáveis sensíveis.
func maskSecret(key, value string) string {
	lower := strings.ToLower(key)
	for _, suffix := range []string{"key", "secret", "token", "password", "senha", "api", "auth", "credential"} {
		if strings.Contains(lower, suffix) {
			if len(value) <= 4 {
				return "****"
			}
			return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
		}
	}
	return value
}
