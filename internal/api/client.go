package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	// normalize: strip trailing slash, ensure /api/v1 suffix
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/api/v1") {
		baseURL = baseURL + "/api/v1"
	}
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func (c *Client) do(method, path string, body any) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := extractErrorMessage(data)
		return nil, &APIError{StatusCode: resp.StatusCode, Message: msg}
	}

	return data, nil
}

func (c *Client) Get(path string) ([]byte, error) {
	return c.do(http.MethodGet, path, nil)
}

func (c *Client) Post(path string, body any) ([]byte, error) {
	return c.do(http.MethodPost, path, body)
}

func (c *Client) Patch(path string, body any) ([]byte, error) {
	return c.do(http.MethodPatch, path, body)
}

func (c *Client) Delete(path string) ([]byte, error) {
	return c.do(http.MethodDelete, path, nil)
}

// Ping valida que o servidor está acessível e o token é válido.
func (c *Client) Ping() error {
	// /health não precisa de token mas confirma que o servidor responde
	req, err := http.NewRequest(http.MethodGet, strings.Replace(c.baseURL, "/api/v1", "/api/health", 1), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("token inválido ou sem permissão")
	}
	return nil
}

func extractErrorMessage(data []byte) string {
	var envelope struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(data, &envelope); err == nil {
		if envelope.Message != "" {
			return envelope.Message
		}
		if envelope.Error != "" {
			return envelope.Error
		}
	}
	if len(data) > 200 {
		return string(data[:200]) + "..."
	}
	return string(data)
}

func decode[T any](data []byte) (T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return v, fmt.Errorf("decoding response: %w", err)
	}
	return v, nil
}
