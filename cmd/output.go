package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

// renderJSON serializes data as JSON when --output json is set.
// Returns true if JSON was rendered (caller should return nil afterward).
// When outputFormat is "table", does nothing and returns false so the caller
// proceeds with normal table rendering.
func renderJSON(v any) bool {
	if outputFormat != "json" {
		return false
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao serializar JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
	return true
}
