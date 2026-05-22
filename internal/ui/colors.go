package ui

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// colorsEnabled controla se escapes ANSI são emitidos.
var colorsEnabled = os.Getenv("NO_COLOR") == "" && os.Getenv("CI") == ""

// useUnicode controla se símbolos Unicode (✓✗→) são usados.
// No Windows, desligamos por padrão devido a problemas de encoding no PowerShell.
var useUnicode = func() bool {
	if os.Getenv("CMX_ASCII") != "" {
		return false
	}
	if runtime.GOOS == "windows" {
		return false
	}
	return true
}()

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	gray   = "\033[90m"
)

func colorize(color, s string) string {
	if !colorsEnabled {
		return s
	}
	return color + s + reset
}

func Bold(s string) string   { return colorize(bold, s) }
func Green(s string) string  { return colorize(green, s) }
func Red(s string) string    { return colorize(red, s) }
func Yellow(s string) string { return colorize(yellow, s) }
func Cyan(s string) string   { return colorize(cyan, s) }
func Gray(s string) string   { return colorize(gray, s) }

// StatusColor formata o status do Coolify com a cor adequada.
func StatusColor(status string) string {
	lower := strings.ToLower(status)
	switch {
	case strings.Contains(lower, "running"):
		return Green(status)
	case strings.Contains(lower, "stopped"), strings.Contains(lower, "exited"):
		return Red(status)
	case strings.Contains(lower, "deploy"), strings.Contains(lower, "starting"), strings.Contains(lower, "restart"):
		return Cyan(status)
	case strings.Contains(lower, "error"), strings.Contains(lower, "failed"):
		return Red(Bold(status))
	case lower == "":
		return Gray("—")
	default:
		return Yellow(status)
	}
}

func Success(msg string) {
	prefix := "OK"
	if useUnicode {
		prefix = "✓"
	}
	if colorsEnabled {
		fmt.Println(Green(prefix) + " " + msg)
	} else {
		fmt.Println(prefix + " " + msg)
	}
}

func Fail(msg string) {
	prefix := "X"
	if useUnicode {
		prefix = "✗"
	}
	if colorsEnabled {
		fmt.Fprintln(os.Stderr, Red(prefix)+" "+msg)
	} else {
		fmt.Fprintln(os.Stderr, prefix+" "+msg)
	}
}

func Info(msg string) {
	prefix := ">"
	if useUnicode {
		prefix = "→"
	}
	if colorsEnabled {
		fmt.Println(Gray(prefix) + " " + msg)
	} else {
		fmt.Println(prefix + " " + msg)
	}
}
