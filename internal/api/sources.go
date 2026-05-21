package api

type GitHubApp struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (c *Client) ListGitHubApps() ([]GitHubApp, error) {
	data, err := c.Get("/github-apps")
	if err != nil {
		return nil, err
	}
	return decode[[]GitHubApp](data)
}
