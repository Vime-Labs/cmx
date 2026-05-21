package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var repoRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)

// Ports valida uma ou mais portas separadas por vírgula (ex: "3000" ou "3000,3001").
func Ports(v string) error {
	if v == "" {
		return fmt.Errorf("porta obrigatória")
	}
	for _, p := range strings.Split(v, ",") {
		p = strings.TrimSpace(p)
		n, err := strconv.Atoi(p)
		if err != nil || n < 1 || n > 65535 {
			return fmt.Errorf("%q não é uma porta válida (1-65535)", p)
		}
	}
	return nil
}

// RepoFormat valida o formato owner/repo.
func RepoFormat(v string) error {
	if !repoRegex.MatchString(v) {
		return fmt.Errorf("formato inválido — use owner/repo (ex: Vime-Labs/meu-app)")
	}
	return nil
}

// BranchName valida que o nome de branch não contém espaços.
func BranchName(v string) error {
	if v == "" {
		return fmt.Errorf("branch obrigatória")
	}
	if strings.ContainsAny(v, " \t") {
		return fmt.Errorf("branch não pode conter espaços")
	}
	return nil
}

// ResourceName valida nomes de recursos (apps, dbs): sem espaços ou barras.
func ResourceName(v string) error {
	if v == "" {
		return fmt.Errorf("nome obrigatório")
	}
	if strings.ContainsAny(v, " \t/\\") {
		return fmt.Errorf("nome não pode conter espaços ou barras")
	}
	return nil
}

// ImageTag valida que a imagem Docker inclui uma tag de versão.
func ImageTag(v string, defaultImg string) error {
	if !strings.Contains(v, ":") {
		return fmt.Errorf("inclua a tag da versão (ex: %s)", defaultImg)
	}
	return nil
}

// KeyValue valida e separa um argumento no formato KEY=VALUE.
// Aceita tanto "KEY=VALUE" quanto dois argumentos separados.
func KeyValue(args []string) (key, value string, err error) {
	switch len(args) {
	case 1:
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("formato inválido — use KEY=VALUE ou KEY VALUE")
		}
		key, value = strings.TrimSpace(parts[0]), parts[1]
	case 2:
		key, value = strings.TrimSpace(args[0]), args[1]
	default:
		return "", "", fmt.Errorf("esperado KEY=VALUE ou KEY VALUE")
	}
	if key == "" {
		return "", "", fmt.Errorf("chave não pode ser vazia")
	}
	return key, value, nil
}
