package api

type Project struct {
	UUID         string        `json:"uuid"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Environments []Environment `json:"environments"`
}

type Environment struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (c *Client) ListProjects() ([]Project, error) {
	data, err := c.Get("/projects")
	if err != nil {
		return nil, err
	}
	return decode[[]Project](data)
}

func (c *Client) GetProject(uuid string) (*Project, error) {
	data, err := c.Get("/projects/" + uuid)
	if err != nil {
		return nil, err
	}
	p, err := decode[Project](data)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Client) CreateProject(name string) (*Project, error) {
	data, err := c.Post("/projects", map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	p, err := decode[Project](data)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
