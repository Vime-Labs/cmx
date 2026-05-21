package ui

import (
	"fmt"
	"regexp"
	"strings"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripANSI remove escape codes para medir a largura real do texto.
func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

// Truncate encurta a string para max runes, adicionando "…" se necessário.
func Truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}

// ShortID retorna os primeiros 8 chars de um UUID para display.
func ShortID(id string) string {
	if len(id) <= 8 {
		return id
	}
	return id[:8] + "…"
}

// Table renderiza tabelas com header bold, separador e colunas alinhadas
// corretamente mesmo com células que contêm códigos ANSI.
type Table struct {
	headers []string
	rows    [][]string
}

func NewTable(headers ...string) *Table {
	return &Table{headers: headers}
}

func (t *Table) AddRow(cells ...string) {
	t.rows = append(t.rows, cells)
}

func (t *Table) Render() {
	n := len(t.headers)

	// calcula larguras usando texto sem ANSI
	widths := make([]int, n)
	for i, h := range t.headers {
		widths[i] = len(h)
	}
	for _, row := range t.rows {
		for i := 0; i < n && i < len(row); i++ {
			if w := len(stripANSI(row[i])); w > widths[i] {
				widths[i] = w
			}
		}
	}

	gap := "   "

	// só imprime header e separador se houver headers não-vazios
	hasHeaders := false
	for _, h := range t.headers {
		if h != "" {
			hasHeaders = true
			break
		}
	}

	if hasHeaders {
		for i, h := range t.headers {
			if i > 0 {
				fmt.Print(gap)
			}
			fmt.Print(Bold(fmt.Sprintf("%-*s", widths[i], h)))
		}
		fmt.Println()

		for i, w := range widths {
			if i > 0 {
				fmt.Print(gap)
			}
			fmt.Print(Gray(strings.Repeat("─", w)))
		}
		fmt.Println()
	}

	// linhas
	for _, row := range t.rows {
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Print(gap)
			}
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			pad := widths[i] - len(stripANSI(cell))
			fmt.Print(cell + strings.Repeat(" ", pad))
		}
		fmt.Println()
	}
}
