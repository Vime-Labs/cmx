// Este arquivo serve como exemplo de padrão de testes usando o mock.
// Ao adicionar novos comandos, escreva testes seguindo esta estrutura.
package mock_test

import (
	"errors"
	"testing"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/api/mock"
)

// TestMockImplementsAPI garante em tempo de compilação que mock.Client
// implementa api.API. Se a interface crescer e o mock não for atualizado,
// este arquivo vai falhar ao compilar antes mesmo de rodar os testes.
var _ api.API = (*mock.Client)(nil)

func TestMockListApps_returnsConfiguredApps(t *testing.T) {
	want := []api.Application{
		{UUID: "abc123", Name: "meu-app", Status: "running"},
	}
	c := &mock.Client{
		ListAppsFunc: func() ([]api.Application, error) {
			return want, nil
		},
	}

	got, err := c.ListApps()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "meu-app" {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestMockListApps_propagatesError(t *testing.T) {
	sentinel := errors.New("api down")
	c := &mock.Client{
		ListAppsFunc: func() ([]api.Application, error) {
			return nil, sentinel
		},
	}

	_, err := c.ListApps()
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestMockDefaultsReturnZeroValues(t *testing.T) {
	// Um mock sem funções configuradas não deve panicar.
	c := &mock.Client{}

	if err := c.Ping(); err != nil {
		t.Errorf("Ping() = %v, want nil", err)
	}
	apps, err := c.ListApps()
	if err != nil || apps != nil {
		t.Errorf("ListApps() = (%v, %v), want (nil, nil)", apps, err)
	}
	msg, err := c.StartApp("any")
	if err != nil || msg != "started" {
		t.Errorf("StartApp() = (%q, %v)", msg, err)
	}
}

func TestMockSetAppEnv_calledWithCorrectArgs(t *testing.T) {
	var gotApp, gotKey, gotVal string
	c := &mock.Client{
		SetAppEnvFunc: func(appID, key, value string) error {
			gotApp, gotKey, gotVal = appID, key, value
			return nil
		},
	}

	if err := c.SetAppEnv("my-app", "NODE_ENV", "production"); err != nil {
		t.Fatal(err)
	}
	if gotApp != "my-app" || gotKey != "NODE_ENV" || gotVal != "production" {
		t.Errorf("SetAppEnv called with (%q, %q, %q)", gotApp, gotKey, gotVal)
	}
}
