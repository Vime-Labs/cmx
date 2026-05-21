package api

type BackupResult struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
	Size    int    `json:"size,omitempty"`
	Status  string `json:"status,omitempty"`
}

func (c *Client) BackupDB(id string) (*BackupResult, error) {
	uuid, err := c.resolveDBUUID(id)
	if err != nil {
		return nil, err
	}
	data, err := c.Post("/databases/"+uuid+"/backup", nil)
	if err != nil {
		return nil, err
	}
	resp, err := decode[BackupResult](data)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
