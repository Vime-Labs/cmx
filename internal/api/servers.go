package api

type Server struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

func (c *Client) ListServers() ([]Server, error) {
	data, err := c.Get("/servers")
	if err != nil {
		return nil, err
	}
	return decode[[]Server](data)
}
