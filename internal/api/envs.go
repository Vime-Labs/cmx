package api

import (
	"fmt"
)

type EnvVar struct {
	UUID      string `json:"uuid"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	IsPreview bool   `json:"is_preview"`
}

func (c *Client) ListAppEnvs(id string) ([]EnvVar, error) {
	uuid, err := c.resolveAppUUID(id)
	if err != nil {
		return nil, err
	}
	data, err := c.Get(fmt.Sprintf("/applications/%s/envs", uuid))
	if err != nil {
		return nil, err
	}
	return decode[[]EnvVar](data)
}

// SetAppEnv cria ou atualiza uma variável de ambiente.
// Verifica primeiro se a chave já existe: se sim, faz PATCH; se não, POST.
func (c *Client) SetAppEnv(appID, key, value string) error {
	uuid, err := c.resolveAppUUID(appID)
	if err != nil {
		return err
	}

	envs, err := c.ListAppEnvs(uuid)
	if err != nil {
		return fmt.Errorf("verificando envs existentes: %w", err)
	}

	body := map[string]string{"key": key, "value": value}
	path := fmt.Sprintf("/applications/%s/envs", uuid)

	for _, e := range envs {
		if e.Key == key {
			_, err = c.Patch(path, body)
			return err
		}
	}

	_, err = c.Post(path, body)
	return err
}

// DeleteAppEnvByKey remove uma variável de ambiente pelo nome da chave.
// Busca o UUID da variável na lista antes de deletar.
func (c *Client) DeleteAppEnvByKey(appID, key string) error {
	uuid, err := c.resolveAppUUID(appID)
	if err != nil {
		return err
	}

	envs, err := c.ListAppEnvs(uuid)
	if err != nil {
		return fmt.Errorf("buscando envs: %w", err)
	}

	for _, e := range envs {
		if e.Key == key {
			_, err = c.Delete(fmt.Sprintf("/applications/%s/envs/%s", uuid, e.UUID))
			return err
		}
	}

	return fmt.Errorf("variável %q não encontrada", key)
}
