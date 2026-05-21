package cmd

import (
	"fmt"
	"strings"

	"github.com/Vime-Labs/cmx/internal/api"
	"github.com/Vime-Labs/cmx/internal/ui"
	"github.com/spf13/cobra"
)

var tunnelUser string

var tunnelCmd = &cobra.Command{
	Use:   "tunnel <app-id|db-id> <porta-local:porta-remota>",
	Short: "Gera comando SSH para túnel com o container",
	Long: `Gera um comando SSH para criar um túnel até o servidor onde o container está rodando.

O comando primeiro tenta resolver o ID como aplicação e, se não encontrar,
tenta como banco de dados. Depois consulta os servidores disponíveis e,
se houver mais de um, pede para você escolher.

Exemplos:
  cmx tunnel meu-app 3306:3306
  cmx tunnel meu-banco 5432:5432
  cmx tunnel --user ubuntu meu-app 3000:3000

O comando SSH será exibido — copie e execute em um terminal separado.`,
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeAppsAndDBs,
	RunE:              runTunnel,
}

func runTunnel(cmd *cobra.Command, args []string) error {
	resourceID := args[0]
	portMapping := args[1]

	// Parse port mapping
	parts := strings.SplitN(portMapping, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("formato inválido. Use <porta-local>:<porta-remota> (ex: 3306:3306)")
	}
	localPort := parts[0]
	remotePort := parts[1]

	client := mustClient()

	// Resolve o recurso para obter o nome e, em seguida, o servidor
	resolvedName, serverIP, err := resolveResourceServer(client, resourceID)
	if err != nil {
		return err
	}

	fmt.Println()
	ui.Info(fmt.Sprintf("Túnel para %s via %s (%s)", resolvedName, serverIP, tunnelUser))
	fmt.Println()
	fmt.Printf("  ssh -L %s:localhost:%s %s@%s\n", localPort, remotePort, tunnelUser, serverIP)
	fmt.Println()
	fmt.Println("  Copie o comando acima e execute em outro terminal.")
	fmt.Println("  Pressione Ctrl+C neste terminal para encerrar quando terminar.")
	fmt.Println()

	// Also show a one-liner with -N (no command execution, just port forward)
	fmt.Println("  Ou, para encerrar automaticamente ao fechar:")
	fmt.Printf("  ssh -L %s:localhost:%s -N %s@%s\n", localPort, remotePort, tunnelUser, serverIP)
	fmt.Println()

	return nil
}

// resolveResourceServer tenta encontrar o servidor onde o recurso está rodando.
// Como a API atual não expõe o servidor de um recurso diretamente,
// consulta a lista de servidores disponíveis e pede para o usuário
// selecionar quando há mais de um.
func resolveResourceServer(client api.API, id string) (string, string, error) {
	var resourceName string

	// Tenta como app primeiro
	app, err := client.GetApp(id)
	if err == nil && app != nil {
		resourceName = app.Name
	}

	// Se não encontrou como app, tenta como DB
	if resourceName == "" {
		db, err := client.GetDB(id)
		if err == nil && db != nil {
			resourceName = db.Name
		}
	}

	if resourceName == "" {
		return "", "", fmt.Errorf("recurso %q não encontrado — verifique o nome ou UUID", id)
	}

	// Lista servidores disponíveis
	servers, err := client.ListServers()
	if err != nil {
		return "", "", fmt.Errorf("falha ao listar servidores: %w", err)
	}

	if len(servers) == 0 {
		return "", "", fmt.Errorf("nenhum servidor disponível — verifique se há servidores configurados")
	}

	if len(servers) == 1 {
		return resourceName, servers[0].IP, nil
	}

	// Múltiplos servidores — pede para o usuário escolher
	names := make([]string, len(servers))
	for i, s := range servers {
		names[i] = fmt.Sprintf("%s (%s)", s.Name, s.IP)
	}
	idx, err := ui.Select(fmt.Sprintf("Em qual servidor o recurso %s está rodando?", resourceName), names)
	if err != nil {
		return "", "", err
	}

	return resourceName, servers[idx].IP, nil
}

func init() {
	tunnelCmd.Flags().StringVar(&tunnelUser, "user", "root", "Usuário SSH do servidor")
	rootCmd.AddCommand(tunnelCmd)
}
