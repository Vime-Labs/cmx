package ui

import (
	"fmt"
	"os"
	"strings"
)

var colorsEnabled = os.Getenv("NO_COLOR") == "" && os.Getenv("CI") == ""

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

func Success(msg string) { fmt.Println(Green("✓") + " " + msg) }
func Fail(msg string)    { fmt.Fprintln(os.Stderr, Red("✗")+" "+msg) }
func Info(msg string)    { fmt.Println(Gray("→") + " " + msg) }
