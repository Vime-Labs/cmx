package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

// statusCategory returns a normalized group label for a resource status.
func statusCategory(s string) string {
	lower := strings.ToLower(s)
	switch {
	case strings.Contains(lower, "running"):
		return "running"
	case strings.Contains(lower, "stopped"), strings.Contains(lower, "exited"):
		return "stopped"
	case strings.Contains(lower, "error"), strings.Contains(lower, "failed"):
		return "error"
	default:
		return "other"
	}
}

// statusCategoryCount returns a count map keyed by category.
func countByStatus[T any](items []T, statusFn func(T) string) map[string]int {
	m := make(map[string]int)
	for _, it := range items {
		m[statusFn(it)]++
	}
	return m
}

const statusBoxWidth = 54

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Exibe um painel consolidado do ambiente",
	Long:  "Mostra um resumo visual com servidores, aplicações, bancos de dados e deployments ativos.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := mustClient()
		spin := ui.NewSpinner("Coletando informações do ambiente")

		// ── chamadas paralelas ──────────────────────────────────────
		var (
			apps    []api.Application
			dbs     []api.Database
			deps    []api.Deployment
			servers []api.Server

			mu   sync.Mutex
			errs []error
			wg   sync.WaitGroup
		)

		collect := func(err error) {
			mu.Lock()
			if err != nil {
				errs = append(errs, err)
			}
			mu.Unlock()
		}

		wg.Add(4)
		go func() {
			defer wg.Done()
			a, err := client.ListApps()
			collect(err)
			mu.Lock()
			apps = a
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			d, err := client.ListDBs()
			collect(err)
			mu.Lock()
			dbs = d
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			d, err := client.ListDeployments()
			collect(err)
			mu.Lock()
			deps = d
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			s, err := client.ListServers()
			collect(err)
			mu.Lock()
			servers = s
			mu.Unlock()
		}()
		wg.Wait()

		if len(errs) > 0 {
			spin.Fail("falhou")
			return errs[0]
		}
		spin.Stop("OK")

		// ── JSON output ─────────────────────────────────────────────
		if renderJSON(struct {
			Servers     []api.Server      `json:"servers"`
			Apps        []api.Application `json:"apps"`
			Databases   []api.Database    `json:"databases"`
			Deployments []api.Deployment  `json:"deployments"`
		}{
			Servers:     servers,
			Apps:        apps,
			Databases:   dbs,
			Deployments: deps,
		}) {
			return nil
		}

		fmt.Println()

		// ── dashboard ──────────────────────────────────────────────
		appCounts := countByStatus(apps, func(a api.Application) string {
			return statusCategory(a.Status)
		})
		dbCounts := countByStatus(dbs, func(d api.Database) string {
			return statusCategory(d.Status)
		})

		totalApps := len(apps)
		totalDBs := len(dbs)

		// Linha superior
		title := " CMX Status "
		fmt.Print("  ")
		fmt.Print(ui.Cyan("┌──" + title + strings.Repeat("─", statusBoxWidth-len(title)-4) + "┐"))
		fmt.Println()

		printLine := func(format string, a ...interface{}) {
			inner := fmt.Sprintf(format, a...)
			fmt.Printf("  │  %-*s  │\n", statusBoxWidth-4, inner)
		}

		// Servidores
		plural := "servidor(es) configurado(s)"
		if len(servers) == 1 {
			plural = "servidor configurado"
		}
		printLine("Servidor:  %s %s", ui.Cyan(fmt.Sprintf("%d", len(servers))), plural)

		if len(servers) > 0 {
			printLine("")
			for _, s := range servers {
				printLine("  %s  %s", ui.Gray(ui.ShortID(s.UUID)), s.Name)
			}
		}
		printLine("")

		// Aplicações
		printLine("Aplicações:   %s total", ui.Bold(fmt.Sprintf("%d", totalApps)))
		printCountLines(printLine, appCounts)

		printLine("")

		// Bancos
		printLine("Bancos:       %s total", ui.Bold(fmt.Sprintf("%d", totalDBs)))
		printCountLines(printLine, dbCounts)

		printLine("")

		// Deployments ativos
		if len(deps) > 0 {
			printLine("Deployments ativos:  %s", ui.Bold(fmt.Sprintf("%d", len(deps))))

			// build app name lookup
			appNames := make(map[string]string, len(apps))
			for _, a := range apps {
				appNames[a.UUID] = a.Name
			}

			for _, d := range deps {
				appName := appNames[d.ApplicationUUID]
				if appName == "" {
					appName = ui.Gray("—")
				}
				line := fmt.Sprintf("  %s  %-16s  %s  %s",
					ui.Gray(ui.ShortID(d.UUID)),
					ui.Truncate(appName, 16),
					ui.StatusColor(ui.Truncate(d.Status, 14)),
					ui.Gray(d.CreatedAt),
				)
				printLine("%s", line)
			}
		} else {
			printLine("Deployments ativos:  %s", ui.Gray("nenhum"))
		}

		// Linha inferior
		fmt.Print("  ")
		fmt.Print(ui.Cyan("└" + strings.Repeat("─", statusBoxWidth-2) + "┘"))
		fmt.Println()
		fmt.Println()

		return nil
	},
}

func printCountLines(printLine func(string, ...interface{}), counts map[string]int) {
	for _, cat := range []string{"running", "stopped", "error", "other"} {
		n := counts[cat]
		if n == 0 {
			continue
		}
		var colorFn func(string) string
		switch cat {
		case "running":
			colorFn = ui.Green
		case "stopped":
			colorFn = ui.Red
		case "error":
			colorFn = func(s string) string { return ui.Red(ui.Bold(s)) }
		default:
			colorFn = ui.Yellow
		}
		tag := ui.Gray("›")
		printLine("  %s %s:  %d", tag, colorFn(cat), n)
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
