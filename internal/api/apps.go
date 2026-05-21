package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Application struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Repository   string `json:"git_repository"`
	Branch       string `json:"git_branch"`
	BuildPack    string `json:"build_pack"`
	PortsExposes string `json:"ports_exposes"`
	Domains      string `json:"fqdn"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type DeployResponse struct {
	Deployments []struct {
		Message        string `json:"message"`
		ResourceUUID   string `json:"resource_uuid"`
		DeploymentUUID string `json:"deployment_uuid"`
	} `json:"deployments"`
}

func (c *Client) ListApps() ([]Application, error) {
	data, err := c.Get("/applications")
	if err != nil {
		return nil, err
	}
	return decode[[]Application](data)
}

func (c *Client) GetApp(id string) (*Application, error) {
	uuid, err := c.resolveAppUUID(id)
	if err != nil {
		return nil, err
	}
	data, err := c.Get("/applications/" + uuid)
	if err != nil {
		return nil, err
	}
	app, err := decode[Application](data)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (c *Client) AppLogs(id string, lines int) (string, error) {
	uuid, err := c.resolveAppUUID(id)
	if err != nil {
		return "", err
	}
	data, err := c.Get(fmt.Sprintf("/applications/%s/logs?lines=%d", uuid, lines))
	if err != nil {
		return "", err
	}
	var resp struct {
		Logs string `json:"logs"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("decoding logs: %w", err)
	}
	return resp.Logs, nil
}

func (c *Client) DeployApp(id string, force bool) (*DeployResponse, error) {
	uuid, err := c.resolveAppUUID(id)
	if err != nil {
		return nil, err
	}
	return c.deployByQuery(fmt.Sprintf("uuid=%s", uuid), force)
}

func (c *Client) DeployByTag(tag string, force bool) (*DeployResponse, error) {
	return c.deployByQuery(fmt.Sprintf("tag=%s", tag), force)
}

func (c *Client) deployByQuery(query string, force bool) (*DeployResponse, error) {
	path := "/deploy?" + query
	if force {
		path += "&force=true"
	}
	data, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	resp, err := decode[DeployResponse](data)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) StartApp(id string) (string, error) {
	return c.appAction(id, "start")
}

func (c *Client) StopApp(id string) (string, error) {
	return c.appAction(id, "stop")
}

func (c *Client) RestartApp(id string) (string, error) {
	return c.appAction(id, "restart")
}

func (c *Client) appAction(id, action string) (string, error) {
	uuid, err := c.resolveAppUUID(id)
	if err != nil {
		return "", err
	}
	data, err := c.Get(fmt.Sprintf("/applications/%s/%s", uuid, action))
	if err != nil {
		return "", err
	}
	var resp struct {
		Message        string `json:"message"`
		DeploymentUUID string `json:"deployment_uuid"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}
	if resp.DeploymentUUID != "" {
		return fmt.Sprintf("%s (deployment: %s)", resp.Message, resp.DeploymentUUID), nil
	}
	return resp.Message, nil
}

type UpdateAppRequest struct {
	Name          *string `json:"name,omitempty"`
	BuildPack     *string `json:"build_pack,omitempty"`
	GitBranch     *string `json:"git_branch,omitempty"`
	GitRepository *string `json:"git_repository,omitempty"`
	PortsExposes  *string `json:"ports_exposes,omitempty"`
	FQDN          *string `json:"fqdn,omitempty"`
	PrivateKeyID  *string `json:"private_key_id,omitempty"`
	DestDir       *string `json:"dest_dir,omitempty"`
}

type CreateAppRequest struct {
	ProjectUUID     string `json:"project_uuid"`
	ServerUUID      string `json:"server_uuid"`
	EnvironmentName string `json:"environment_name"`
	GitHubAppUUID   string `json:"github_app_uuid"`
	GitRepository   string `json:"git_repository"`
	GitBranch       string `json:"git_branch"`
	BuildPack       string `json:"build_pack"`
	PortsExposes    string `json:"ports_exposes"`
	Name            string `json:"name,omitempty"`
	FQDN            string `json:"fqdn,omitempty"`
}

type CreateAppResponse struct {
	UUID           string `json:"uuid"`
	DeploymentUUID string `json:"deployment_uuid"`
}

func (c *Client) CreateApp(req CreateAppRequest) (*CreateAppResponse, error) {
	data, err := c.Post("/applications/private-github-app", req)
	if err != nil {
		return nil, err
	}
	resp, err := decode[CreateAppResponse](data)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// resolveAppUUID aceita UUID direto ou nome (busca na lista).
func (c *Client) resolveAppUUID(id string) (string, error) {
	if looksLikeUUID(id) {
		return id, nil
	}
	apps, err := c.ListApps()
	if err != nil {
		return "", err
	}
	var matches []Application
	for _, a := range apps {
		if strings.EqualFold(a.Name, id) {
			matches = append(matches, a)
		}
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("aplicação %q não encontrada", id)
	case 1:
		return matches[0].UUID, nil
	default:
		names := make([]string, len(matches))
		for i, m := range matches {
			names[i] = fmt.Sprintf("%s (%s)", m.Name, m.UUID)
		}
		return "", fmt.Errorf("nome ambíguo, use o UUID:\n  %s", strings.Join(names, "\n  "))
	}
}

func looksLikeUUID(s string) bool {
	// UUIDs do Coolify têm 8 chars hex sem traço, ex: "a1b2c3d4"
	if len(s) < 8 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') || c == '-') {
			return false
		}
	}
	return true
}
