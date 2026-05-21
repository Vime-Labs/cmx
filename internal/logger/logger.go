package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Entry representa uma única linha do histórico de atividade.
type Entry struct {
	Timestamp    string `json:"timestamp"`
	Command      string `json:"command"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	Detail       string `json:"detail,omitempty"`
	Status       string `json:"status"`
	Duration     string `json:"duration,omitempty"`
}

// Action é uma constante auxiliar para os comandos.
type Action string

const (
	ActionAppCreate    Action = "apps create"
	ActionAppDeploy    Action = "apps deploy"
	ActionAppStart     Action = "apps start"
	ActionAppStop      Action = "apps stop"
	ActionAppRestart   Action = "apps restart"
	ActionEnvSet       Action = "apps envs set"
	ActionEnvDelete    Action = "apps envs delete"
	ActionDomainAdd    Action = "domain add"
	ActionDomainRemove Action = "domain remove"
	ActionDBCreate     Action = "dbs create"
	ActionDBStart      Action = "dbs start"
	ActionDBStop       Action = "dbs stop"
	ActionDBRestart    Action = "dbs restart"
	ActionDBBackup     Action = "dbs backup"
	ActionDeploy       Action = "deploy"
	ActionDeployCancel Action = "deployments cancel"
	ActionConfigure    Action = "configure"
)

const ResourceApp = "app"
const ResourceDB = "database"
const ResourceEnv = "env"
const ResourceCfg = "config"
const ResourceDepl = "deployment"

// logPath retorna o caminho do arquivo de log.
func logPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cmx", "activity.log"), nil
}

// Log registra uma entrada no histórico de atividade.
// Se a operação foi bem-sucedida, passe status="success"; caso contrário, "error".
// O duration é opcional (passe "" se não quiser registrar tempo).
func Log(action Action, resourceType, resourceName, detail, status string, duration ...time.Duration) {
	entry := Entry{
		Timestamp:    time.Now().Format(time.RFC3339),
		Command:      string(action),
		ResourceType: resourceType,
		ResourceName: resourceName,
		Detail:       detail,
		Status:       status,
	}
	if len(duration) > 0 {
		entry.Duration = fmtDuration(duration[0])
	}

	path, err := logPath()
	if err != nil {
		return // silent fail — log não deve quebrar o fluxo
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	line, _ := json.Marshal(entry)
	f.Write(line)
	f.Write([]byte("\n"))
}

// List retorna todas as entradas de log, da mais recente para a mais antiga.
// Se max > 0, limita o retorno às N mais recentes.
// Se filterAction não for vazio, filtra pelo comando.
func List(max int, filterAction string) ([]Entry, error) {
	path, err := logPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading activity log: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	// reverse: mais recente primeiro
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	var entries []Entry
	for _, line := range lines {
		if line == "" {
			continue
		}
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}
		if filterAction != "" && !strings.EqualFold(e.Command, filterAction) &&
			!strings.HasPrefix(strings.ToLower(e.Command), strings.ToLower(filterAction)) {
			continue
		}
		entries = append(entries, e)
		if max > 0 && len(entries) >= max {
			break
		}
	}

	// Estável: se não pediu limite e não filtrou, mantém ordem cronológica reversa
	// (já está em ordem reversa)
	return entries, nil
}

// Clear apaga todo o histórico de atividade.
func Clear() error {
	path, err := logPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}

// ── helpers ─────────────────────────────────────────────────────────────────

func fmtDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		sec := d.Seconds()
		return fmt.Sprintf("%.1fs", sec)
	}
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm%ds", m, s)
}

// ── Stats ───────────────────────────────────────────────────────────────────

type Stats struct {
	Total      int            `json:"total"`
	Success    int            `json:"success"`
	Errors     int            `json:"errors"`
	ByCommand  map[string]int `json:"by_command"`
	ByResource map[string]int `json:"by_resource"`
}

// ComputeStats percorre todas as entradas e calcula estatísticas agregadas.
func ComputeStats() (*Stats, error) {
	entries, err := List(0, "")
	if err != nil {
		return nil, err
	}

	s := &Stats{
		ByCommand:  make(map[string]int),
		ByResource: make(map[string]int),
	}
	for _, e := range entries {
		s.Total++
		if e.Status == "success" {
			s.Success++
		} else {
			s.Errors++
		}
		s.ByCommand[e.Command]++
		s.ByResource[e.ResourceType]++
	}
	return s, nil
}

// GroupByDay agrupa as entradas por dia (YYYY-MM-DD) em ordem cronológica.
func GroupByDay(entries []Entry) map[string][]Entry {
	groups := make(map[string][]Entry)
	for _, e := range entries {
		day := e.Timestamp[:10] // "2025-01-15"
		groups[day] = append(groups[day], e)
	}
	return groups
}

// SortedDays retorna as chaves de GroupByDay ordenadas (mais recente primeiro).
func SortedDays(groups map[string][]Entry) []string {
	days := make([]string, 0, len(groups))
	for d := range groups {
		days = append(days, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(days)))
	return days
}
