package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Database struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Type        string `json:"type"` // ex: "standalone-postgresql"
	Image       string `json:"image"`
	IsPublic    bool   `json:"is_public"`
	PublicPort  int    `json:"public_port"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	// Credenciais (retornadas pela API no GET detalhado)
	InternalDBURL string `json:"internal_db_url,omitempty"`
	Password      string `json:"password,omitempty"`
	Username      string `json:"username,omitempty"`
	DefaultDB     string `json:"default_database,omitempty"`
}

// DisplayType normaliza o tipo retornado pela API para exibição.
// A API retorna "standalone-postgresql", "standalone-redis", etc.
// Se vazio, deriva da imagem Docker.
func (d *Database) DisplayType() string {
	t := strings.TrimPrefix(d.Type, "standalone-")
	if t != "" && t != d.Type {
		return t
	}
	// fallback: deriva da imagem
	img := strings.ToLower(d.Image)
	for _, candidate := range []string{"postgres", "mysql", "mariadb", "mongo", "redis", "dragonfly", "keydb", "clickhouse"} {
		if strings.HasPrefix(img, candidate) {
			return candidate
		}
	}
	if t != "" {
		return t
	}
	return "—"
}

func (c *Client) ListDBs() ([]Database, error) {
	data, err := c.Get("/databases")
	if err != nil {
		return nil, err
	}
	return decode[[]Database](data)
}

func (c *Client) GetDB(id string) (*Database, error) {
	uuid, err := c.resolveDBUUID(id)
	if err != nil {
		return nil, err
	}
	data, err := c.Get("/databases/" + uuid)
	if err != nil {
		return nil, err
	}
	db, err := decode[Database](data)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (c *Client) StartDB(id string) (string, error) {
	return c.dbAction(id, "start")
}

func (c *Client) StopDB(id string) (string, error) {
	return c.dbAction(id, "stop")
}

func (c *Client) RestartDB(id string) (string, error) {
	return c.dbAction(id, "restart")
}

func (c *Client) dbAction(id, action string) (string, error) {
	uuid, err := c.resolveDBUUID(id)
	if err != nil {
		return "", err
	}
	data, err := c.Get(fmt.Sprintf("/databases/%s/%s", uuid, action))
	if err != nil {
		return "", err
	}
	var resp struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}
	return resp.Message, nil
}

type CreateDBRequest struct {
	ProjectUUID     string `json:"project_uuid"`
	ServerUUID      string `json:"server_uuid"`
	EnvironmentName string `json:"environment_name"`
	Name            string `json:"name,omitempty"`
	Image           string `json:"image,omitempty"`
	Password        string `json:"password,omitempty"`
	IsPublic        bool   `json:"is_public,omitempty"`
	PublicPort      int    `json:"public_port,omitempty"`
}

type UpdateDBRequest struct {
	Name        *string `json:"name,omitempty"`
	Image       *string `json:"image,omitempty"`
	IsPublic    *bool   `json:"is_public,omitempty"`
	PublicPort  *int    `json:"public_port,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CreateDBResponse struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
}

var DBDefaultImages = map[string]string{
	"postgresql": "postgres:16",
	"mysql":      "mysql:8",
	"mariadb":    "mariadb:11",
	"mongodb":    "mongo:7",
	"redis":      "redis:7",
	"dragonfly":  "docker.dragonflydb.io/dragonflydb/dragonfly:latest",
	"keydb":      "eqalpha/keydb:latest",
	"clickhouse": "clickhouse/clickhouse-server:24",
}

func (c *Client) CreateDB(dbType string, req CreateDBRequest) (*CreateDBResponse, error) {
	data, err := c.Post("/databases/"+dbType, req)
	if err != nil {
		return nil, err
	}
	resp, err := decode[CreateDBResponse](data)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) resolveDBUUID(id string) (string, error) {
	if looksLikeUUID(id) {
		return id, nil
	}
	dbs, err := c.ListDBs()
	if err != nil {
		return "", err
	}
	var matches []Database
	for _, d := range dbs {
		if strings.EqualFold(d.Name, id) {
			matches = append(matches, d)
		}
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("banco %q não encontrado", id)
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
