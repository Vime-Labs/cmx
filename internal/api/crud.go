package api

// ── App CRUD ────────────────────────────────────────────────────────────────

func (c *Client) DeleteApp(id string) error {
	uuid, err := c.resolveAppUUID(id)
	if err != nil {
		return err
	}
	_, err = c.Delete("/applications/" + uuid)
	return err
}

func (c *Client) UpdateApp(id string, req UpdateAppRequest) (*Application, error) {
	uuid, err := c.resolveAppUUID(id)
	if err != nil {
		return nil, err
	}
	data, err := c.Patch("/applications/"+uuid, req)
	if err != nil {
		return nil, err
	}
	app, err := decode[Application](data)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// ── DB CRUD ─────────────────────────────────────────────────────────────────

func (c *Client) DeleteDB(id string) error {
	uuid, err := c.resolveDBUUID(id)
	if err != nil {
		return err
	}
	_, err = c.Delete("/databases/" + uuid)
	return err
}

func (c *Client) UpdateDB(id string, req UpdateDBRequest) (*Database, error) {
	uuid, err := c.resolveDBUUID(id)
	if err != nil {
		return nil, err
	}
	data, err := c.Patch("/databases/"+uuid, req)
	if err != nil {
		return nil, err
	}
	db, err := decode[Database](data)
	if err != nil {
		return nil, err
	}
	return &db, nil
}
