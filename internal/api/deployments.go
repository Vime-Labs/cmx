package api

import "fmt"

type Deployment struct {
	UUID            string `json:"uuid"`
	Status          string `json:"status"`
	ApplicationUUID string `json:"application_uuid"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func (c *Client) ListDeployments() ([]Deployment, error) {
	data, err := c.Get("/deployments")
	if err != nil {
		return nil, err
	}
	return decode[[]Deployment](data)
}

func (c *Client) GetDeployment(uuid string) (*Deployment, error) {
	data, err := c.Get("/deployments/" + uuid)
	if err != nil {
		return nil, err
	}
	d, err := decode[Deployment](data)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (c *Client) CancelDeployment(uuid string) error {
	_, err := c.Post("/deployments/"+uuid+"/cancel", nil)
	return err
}

func (c *Client) ListAppDeployments(appID string, skip, take int) ([]Deployment, error) {
	uuid, err := c.resolveAppUUID(appID)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/deployments/applications/%s?skip=%d&take=%d", uuid, skip, take)
	data, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	return decode[[]Deployment](data)
}
