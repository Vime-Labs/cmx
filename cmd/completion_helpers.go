package cmd

import (
	"fmt"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/config"
	"github.com/spf13/cobra"
)

// completeApps completes app names (not UUIDs) for shell autocomplete.
// Uses a client with short timeout to avoid blocking the shell.
func completeApps(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg, err := tryLoadConfig()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	client := api.NewClient(cfg.URL, cfg.Token)
	apps, err := client.ListApps()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var names []string
	for _, a := range apps {
		names = append(names, fmt.Sprintf("%s\t%s - %s", a.Name, a.Status, a.Repository))
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

// completeDBs completes database names for shell autocomplete.
func completeDBs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg, err := tryLoadConfig()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	client := api.NewClient(cfg.URL, cfg.Token)
	dbs, err := client.ListDBs()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var names []string
	for _, d := range dbs {
		names = append(names, fmt.Sprintf("%s\t%s - %s", d.Name, d.Status, d.DisplayType()))
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

// completeAppsAndDBs tries app names first, then DB names.
func completeAppsAndDBs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Try both and merge
	apps, _ := completeApps(cmd, args, toComplete)
	dbs, _ := completeDBs(cmd, args, toComplete)
	all := append(apps, dbs...)
	if len(all) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return all, cobra.ShellCompDirectiveNoFileComp
}

// tryLoadConfig loads config without exiting on error (for completion).
func tryLoadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}
