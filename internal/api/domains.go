package api

type Domain struct {
	UUID      string `json:"uuid"`
	Domain    string `json:"domain"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *Client) ListDomains(appID string) ([]Domain, error) {
	uuid, err := c.resolveAppUUID(appID)
	if err != nil {
		return nil, err
	}
	data, err := c.Get("/applications/" + uuid + "/domains")
	if err != nil {
		return nil, err
	}
	return decode[[]Domain](data)
}

func (c *Client) AddDomain(appID, domain string) (*Domain, error) {
	uuid, err := c.resolveAppUUID(appID)
	if err != nil {
		return nil, err
	}
	data, err := c.Post("/applications/"+uuid+"/domains", map[string]string{"domain": domain})
	if err != nil {
		return nil, err
	}
	d, err := decode[Domain](data)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (c *Client) RemoveDomain(appID, domainUUID string) error {
	uuid, err := c.resolveAppUUID(appID)
	if err != nil {
		return err
	}
	_, err = c.Delete("/applications/" + uuid + "/domains/" + domainUUID)
	return err
}
