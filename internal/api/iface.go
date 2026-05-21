package api

// API é a interface que abstrai todas as operações do Coolify.
// Comandos dependem desta interface, não do client concreto,
// permitindo mocks em testes sem nenhuma chamada HTTP real.
type API interface {
	Ping() error

	// Aplicações
	ListApps() ([]Application, error)
	GetApp(id string) (*Application, error)
	AppLogs(id string, lines int) (string, error)
	DeployApp(id string, force bool) (*DeployResponse, error)
	DeployByTag(tag string, force bool) (*DeployResponse, error)
	StartApp(id string) (string, error)
	StopApp(id string) (string, error)
	RestartApp(id string) (string, error)
	CreateApp(req CreateAppRequest) (*CreateAppResponse, error)

	// Variáveis de ambiente
	ListAppEnvs(id string) ([]EnvVar, error)
	SetAppEnv(appID, key, value string) error
	DeleteAppEnvByKey(appID, key string) error

	// Bancos de dados
	ListDBs() ([]Database, error)
	GetDB(id string) (*Database, error)
	StartDB(id string) (string, error)
	StopDB(id string) (string, error)
	RestartDB(id string) (string, error)
	CreateDB(dbType string, req CreateDBRequest) (*CreateDBResponse, error)

	// Projetos e infraestrutura
	ListProjects() ([]Project, error)
	GetProject(uuid string) (*Project, error)
	ListServers() ([]Server, error)
	ListGitHubApps() ([]GitHubApp, error)

	// Deployments
	ListDeployments() ([]Deployment, error)
	GetDeployment(uuid string) (*Deployment, error)
	CancelDeployment(uuid string) error
	ListAppDeployments(appID string, skip, take int) ([]Deployment, error)
}
