// Package mock fornece um cliente Coolify falso para testes.
// Use os campos Func para injetar o comportamento desejado em cada teste.
//
// Exemplo:
//
//	c := &mock.Client{
//	    ListAppsFunc: func() ([]api.Application, error) {
//	        return []api.Application{{UUID: "abc", Name: "meu-app"}}, nil
//	    },
//	}
package mock

import "github.com/Vime-Labs/cmx/internal/api"

// Client implementa api.API com funções configuráveis por teste.
// Qualquer método não configurado retorna zero value + nil error.
type Client struct {
	PingFunc func() error

	ListAppsFunc    func() ([]api.Application, error)
	GetAppFunc      func(id string) (*api.Application, error)
	AppLogsFunc     func(id string, lines int) (string, error)
	DeployAppFunc   func(id string, force bool) (*api.DeployResponse, error)
	DeployByTagFunc func(tag string, force bool) (*api.DeployResponse, error)
	StartAppFunc    func(id string) (string, error)
	StopAppFunc     func(id string) (string, error)
	RestartAppFunc  func(id string) (string, error)
	CreateAppFunc   func(req api.CreateAppRequest) (*api.CreateAppResponse, error)

	ListAppEnvsFunc       func(id string) ([]api.EnvVar, error)
	SetAppEnvFunc         func(appID, key, value string) error
	DeleteAppEnvByKeyFunc func(appID, key string) error

	ListDBsFunc   func() ([]api.Database, error)
	GetDBFunc     func(id string) (*api.Database, error)
	StartDBFunc   func(id string) (string, error)
	StopDBFunc    func(id string) (string, error)
	RestartDBFunc func(id string) (string, error)
	CreateDBFunc  func(dbType string, req api.CreateDBRequest) (*api.CreateDBResponse, error)

	ListProjectsFunc   func() ([]api.Project, error)
	GetProjectFunc     func(uuid string) (*api.Project, error)
	ListServersFunc    func() ([]api.Server, error)
	ListGitHubAppsFunc func() ([]api.GitHubApp, error)

	ListDomainsFunc  func(appID string) ([]api.Domain, error)
	AddDomainFunc    func(appID, domain string) (*api.Domain, error)
	RemoveDomainFunc func(appID, domainUUID string) error

	BackupDBFunc func(id string) (*api.BackupResult, error)

	ListDeploymentsFunc    func() ([]api.Deployment, error)
	GetDeploymentFunc      func(uuid string) (*api.Deployment, error)
	CancelDeploymentFunc   func(uuid string) error
	ListAppDeploymentsFunc func(appID string, skip, take int) ([]api.Deployment, error)
}

func (m *Client) Ping() error {
	if m.PingFunc != nil {
		return m.PingFunc()
	}
	return nil
}

func (m *Client) ListApps() ([]api.Application, error) {
	if m.ListAppsFunc != nil {
		return m.ListAppsFunc()
	}
	return nil, nil
}

func (m *Client) GetApp(id string) (*api.Application, error) {
	if m.GetAppFunc != nil {
		return m.GetAppFunc(id)
	}
	return &api.Application{}, nil
}

func (m *Client) AppLogs(id string, lines int) (string, error) {
	if m.AppLogsFunc != nil {
		return m.AppLogsFunc(id, lines)
	}
	return "", nil
}

func (m *Client) DeployApp(id string, force bool) (*api.DeployResponse, error) {
	if m.DeployAppFunc != nil {
		return m.DeployAppFunc(id, force)
	}
	return &api.DeployResponse{}, nil
}

func (m *Client) DeployByTag(tag string, force bool) (*api.DeployResponse, error) {
	if m.DeployByTagFunc != nil {
		return m.DeployByTagFunc(tag, force)
	}
	return &api.DeployResponse{}, nil
}

func (m *Client) StartApp(id string) (string, error) {
	if m.StartAppFunc != nil {
		return m.StartAppFunc(id)
	}
	return "started", nil
}

func (m *Client) StopApp(id string) (string, error) {
	if m.StopAppFunc != nil {
		return m.StopAppFunc(id)
	}
	return "stopped", nil
}

func (m *Client) RestartApp(id string) (string, error) {
	if m.RestartAppFunc != nil {
		return m.RestartAppFunc(id)
	}
	return "restarted", nil
}

func (m *Client) CreateApp(req api.CreateAppRequest) (*api.CreateAppResponse, error) {
	if m.CreateAppFunc != nil {
		return m.CreateAppFunc(req)
	}
	return &api.CreateAppResponse{}, nil
}

func (m *Client) ListAppEnvs(id string) ([]api.EnvVar, error) {
	if m.ListAppEnvsFunc != nil {
		return m.ListAppEnvsFunc(id)
	}
	return nil, nil
}

func (m *Client) SetAppEnv(appID, key, value string) error {
	if m.SetAppEnvFunc != nil {
		return m.SetAppEnvFunc(appID, key, value)
	}
	return nil
}

func (m *Client) DeleteAppEnvByKey(appID, key string) error {
	if m.DeleteAppEnvByKeyFunc != nil {
		return m.DeleteAppEnvByKeyFunc(appID, key)
	}
	return nil
}

func (m *Client) ListDBs() ([]api.Database, error) {
	if m.ListDBsFunc != nil {
		return m.ListDBsFunc()
	}
	return nil, nil
}

func (m *Client) GetDB(id string) (*api.Database, error) {
	if m.GetDBFunc != nil {
		return m.GetDBFunc(id)
	}
	return &api.Database{}, nil
}

func (m *Client) StartDB(id string) (string, error) {
	if m.StartDBFunc != nil {
		return m.StartDBFunc(id)
	}
	return "started", nil
}

func (m *Client) StopDB(id string) (string, error) {
	if m.StopDBFunc != nil {
		return m.StopDBFunc(id)
	}
	return "stopped", nil
}

func (m *Client) RestartDB(id string) (string, error) {
	if m.RestartDBFunc != nil {
		return m.RestartDBFunc(id)
	}
	return "restarted", nil
}

func (m *Client) CreateDB(dbType string, req api.CreateDBRequest) (*api.CreateDBResponse, error) {
	if m.CreateDBFunc != nil {
		return m.CreateDBFunc(dbType, req)
	}
	return &api.CreateDBResponse{}, nil
}

func (m *Client) ListProjects() ([]api.Project, error) {
	if m.ListProjectsFunc != nil {
		return m.ListProjectsFunc()
	}
	return nil, nil
}

func (m *Client) GetProject(uuid string) (*api.Project, error) {
	if m.GetProjectFunc != nil {
		return m.GetProjectFunc(uuid)
	}
	return &api.Project{}, nil
}

func (m *Client) ListServers() ([]api.Server, error) {
	if m.ListServersFunc != nil {
		return m.ListServersFunc()
	}
	return nil, nil
}

func (m *Client) ListGitHubApps() ([]api.GitHubApp, error) {
	if m.ListGitHubAppsFunc != nil {
		return m.ListGitHubAppsFunc()
	}
	return nil, nil
}

func (m *Client) ListDomains(appID string) ([]api.Domain, error) {
	if m.ListDomainsFunc != nil {
		return m.ListDomainsFunc(appID)
	}
	return nil, nil
}

func (m *Client) AddDomain(appID, domain string) (*api.Domain, error) {
	if m.AddDomainFunc != nil {
		return m.AddDomainFunc(appID, domain)
	}
	return &api.Domain{}, nil
}

func (m *Client) RemoveDomain(appID, domainUUID string) error {
	if m.RemoveDomainFunc != nil {
		return m.RemoveDomainFunc(appID, domainUUID)
	}
	return nil
}

func (m *Client) BackupDB(id string) (*api.BackupResult, error) {
	if m.BackupDBFunc != nil {
		return m.BackupDBFunc(id)
	}
	return &api.BackupResult{}, nil
}

func (m *Client) ListDeployments() ([]api.Deployment, error) {
	if m.ListDeploymentsFunc != nil {
		return m.ListDeploymentsFunc()
	}
	return nil, nil
}

func (m *Client) GetDeployment(uuid string) (*api.Deployment, error) {
	if m.GetDeploymentFunc != nil {
		return m.GetDeploymentFunc(uuid)
	}
	return &api.Deployment{}, nil
}

func (m *Client) CancelDeployment(uuid string) error {
	if m.CancelDeploymentFunc != nil {
		return m.CancelDeploymentFunc(uuid)
	}
	return nil
}

func (m *Client) ListAppDeployments(appID string, skip, take int) ([]api.Deployment, error) {
	if m.ListAppDeploymentsFunc != nil {
		return m.ListAppDeploymentsFunc(appID, skip, take)
	}
	return nil, nil
}
