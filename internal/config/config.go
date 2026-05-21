package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cmx", "config.yaml"), nil
}

func Load() (*Config, error) {
	cfg := &Config{
		URL:   os.Getenv("CMX_URL"),
		Token: os.Getenv("CMX_TOKEN"),
	}

	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err == nil {
		var fileCfg Config
		if err := yaml.Unmarshal(data, &fileCfg); err != nil {
			return nil, fmt.Errorf("parsing config: %w", err)
		}
		// env vars override file
		if cfg.URL == "" {
			cfg.URL = fileCfg.URL
		}
		if cfg.Token == "" {
			cfg.Token = fileCfg.Token
		}
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func (c *Config) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("Coolify URL not set. Run `cmx configure` or set CMX_URL")
	}
	if c.Token == "" {
		return fmt.Errorf("API token not set. Run `cmx configure` or set CMX_TOKEN")
	}
	return nil
}
