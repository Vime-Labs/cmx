package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func readLine() string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// Select exibe uma lista numerada e retorna o índice escolhido.
// Auto-seleciona se só houver uma opção.
func Select(label string, options []string) (int, error) {
	if len(options) == 0 {
		return 0, fmt.Errorf("nenhuma opção disponível para %q", label)
	}
	if len(options) == 1 {
		fmt.Printf("%s: %s (único disponível)\n", label, options[0])
		return 0, nil
	}
	fmt.Printf("\n%s:\n", label)
	for i, o := range options {
		fmt.Printf("  [%d] %s\n", i+1, o)
	}
	for {
		fmt.Print("Escolha: ")
		input := readLine()
		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > len(options) {
			fmt.Printf("  ✗ Digite um número entre 1 e %d.\n", len(options))
			continue
		}
		return n - 1, nil
	}
}

// Input pede texto obrigatório com default opcional e validador.
func Input(label, defaultVal string, validate func(string) error) (string, error) {
	hint := ""
	if defaultVal != "" {
		hint = fmt.Sprintf(" [%s]", defaultVal)
	}
	for {
		fmt.Printf("%s%s: ", label, hint)
		val := readLine()
		if val == "" {
			val = defaultVal
		}
		if val == "" {
			fmt.Println("  ✗ Campo obrigatório.")
			continue
		}
		if validate != nil {
			if err := validate(val); err != nil {
				fmt.Printf("  ✗ %v\n", err)
				continue
			}
		}
		return val, nil
	}
}

// InputOptional pede texto opcional (vazio é permitido).
func InputOptional(label string) string {
	fmt.Printf("%s (opcional): ", label)
	return readLine()
}

// Confirm pergunta sim/não com default não.
func Confirm(label string) bool {
	fmt.Printf("\n%s [s/N]: ", label)
	val := strings.ToLower(readLine())
	return val == "s" || val == "sim" || val == "y" || val == "yes"
}

// Summary exibe um bloco de resumo formatado antes de confirmar.
func Summary(fields [][2]string) {
	fmt.Println()
	fmt.Println("  ┌─ Resumo ──────────────────────────────┐")
	for _, f := range fields {
		label := f[0]
		value := f[1]
		if value == "" {
			value = "—"
		}
		fmt.Printf("  │  %-16s %s\n", label, value)
	}
	fmt.Println("  └───────────────────────────────────────┘")
}
