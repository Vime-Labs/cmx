# cmx

CLI da Vime Labs para gerenciar aplicações e bancos de dados nos servidores remotos via Coolify.

## Instalação

> **Pré-requisito:** [GitHub CLI](https://cli.github.com) autenticado (`gh auth login`).

**macOS / Linux**
```sh
gh release download --repo Vime-Labs/cmx \
  --pattern "cmx-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')" \
  --output /tmp/cmx --clobber \
  && chmod +x /tmp/cmx \
  && sudo mv /tmp/cmx /usr/local/bin/cmx
```

**Windows (PowerShell)** — copie a linha inteira abaixo e cole no PowerShell:
```powershell
$dir="$env:USERPROFILE\.cmx\bin"; New-Item -Force -ItemType Directory $dir | Out-Null; gh release download --repo Vime-Labs/cmx --pattern cmx-windows-amd64.exe --output "$dir\cmx.exe" --clobber; $p=[Environment]::GetEnvironmentVariable("PATH","User"); if ($p -notlike "*$dir*") { [Environment]::SetEnvironmentVariable("PATH","$p;$dir","User") }; $env:PATH+=";$dir"
```

Após a instalação no Windows, reinicie o terminal para o PATH atualizar.

## Configuração

Na primeira vez, aponte o CMX para o seu Coolify:

```sh
cmx configure
```

Você será solicitado pela URL do servidor e pelo token de API.  
O token é gerado em **Keys & Tokens → API tokens** no painel do Coolify.

Alternativamente, via variáveis de ambiente (útil em CI):

```sh
export CMX_URL=http://seu-servidor:8000
export CMX_TOKEN=seu-token
```

## Uso

### Projetos

```sh
cmx projects list                  # lista todos os projetos
cmx projects get <uuid>            # detalhes de um projeto
cmx projects create <nome>         # cria um novo projeto
```

### Sources / GitHub Apps

```sh
cmx sources list                   # lista GitHub Apps configurados
cmx sources list -o json           # em JSON (UUID + nome)
```

Use o UUID ou nome listado aqui no `--github-app` do `cmx apps create`.

### Aplicações

```sh
cmx apps list                        # lista todas as aplicações
cmx apps get <uuid|nome>             # detalhes de uma aplicação
cmx apps create                      # wizard interativo (cria projeto se necessário)
cmx apps create -p <projeto> -e prod -s <servidor> -g <gh-app> -r owner/repo  # não-interativo
cmx apps update <uuid|nome> --name "novo-nome"  # atualiza configurações
cmx apps delete <uuid|nome>          # remove (com confirmação)
cmx apps delete <uuid|nome> --yes    # remove sem confirmação
cmx apps deploy <uuid|nome>          # dispara deploy
cmx apps deploy <uuid|nome> --force  # deploy sem cache
cmx apps deploy <uuid|nome> --wait   # deploy bloqueante (acompanha até concluir)
cmx apps logs <uuid|nome>            # últimas 100 linhas de log
cmx apps logs <uuid|nome> -n 200     # N linhas de log
cmx apps logs <uuid|nome> --follow   # acompanha logs em tempo real
cmx apps start <uuid|nome>
cmx apps stop <uuid|nome>
cmx apps restart <uuid|nome>
```

### Variáveis de ambiente

```sh
cmx apps envs list <uuid|nome>
cmx apps envs set <uuid|nome> CHAVE=VALOR
cmx apps envs set <uuid|nome> CHAVE VALOR
cmx apps envs delete <uuid|nome> CHAVE
```

### Bancos de dados

```sh
cmx dbs list
cmx dbs get <uuid|nome>
cmx dbs create                       # wizard interativo
cmx dbs create -p <projeto> -e prod -s <servidor> -t postgresql -n meu-db  # não-interativo
cmx dbs update <uuid|nome> --name "novo-nome"  # atualiza configurações
cmx dbs delete <uuid|nome>          # remove (com confirmação)
cmx dbs delete <uuid|nome> --yes    # remove sem confirmação
cmx dbs start <uuid|nome>
cmx dbs stop <uuid|nome>
cmx dbs restart <uuid|nome>
cmx dbs backup <uuid|nome>           # dispara backup
```

### Domínios

```sh
cmx domain list <app-id>             # lista domínios de uma app
cmx domain add <app-id> <domínio>    # adiciona domínio
cmx domain remove <app-id> <uuid>    # remove domínio
```

### Deployments

```sh
cmx deploy <uuid|nome>               # deploy por nome ou UUID
cmx deploy <uuid|nome> --force       # deploy sem cache
cmx deploy <uuid|nome> --wait        # deploy bloqueante (acompanha até concluir)
cmx deploy <uuid|nome> --wait --check  # deploy com gating CI/CD (exit code)
cmx deploy --tag production          # deploy de todos com a tag
cmx deploy --tag production --wait   # deploy bloqueante por tag

cmx deployments list                 # deployments em andamento
cmx deployments get <uuid>           # detalhes de um deployment
cmx deployments history <uuid|nome>  # histórico da aplicação
cmx deployments cancel <uuid>        # cancela deployment ativo
```

### Configuração

```sh
cmx configure                        # wizard interativo
cmx configure --url <url> --token <token>  # não-interativo (para scripts/IA)
cmx config get                       # exibe config atual
cmx config get url                   # exibe apenas a URL
cmx config get token                 # exibe apenas o token
cmx config set url <url>             # atualiza a URL
cmx config set token <token>         # atualiza o token
```

### Túnel SSH

```sh
cmx tunnel <app|db> <local:remoto>   # gera comando SSH para túnel
cmx tunnel meu-app 3306:3306
cmx tunnel --user ubuntu meu-app 3000:3000
```

### Status / Dashboard

```sh
cmx status                           # visão geral do ambiente
```

### Saída JSON

Todos os comandos `list` e `get` aceitam `--output json` (ou `-o json`)
para emitir JSON puro — ideal para scripts e automação:

```sh
cmx apps list -o json
cmx dbs get meu-banco -o json
cmx deployments list --output json
cmx status -o json
cmx log --output json
```

### Histórico de atividade

```sh
cmx log                              # últimas 20 ações
cmx log -n 50                        # últimas 50 ações
cmx log --cmd deploy                 # filtrar por comando
cmx log --clear                      # apagar o histórico
cmx log --clear --yes                # apagar sem confirmação
```

Toda ação executada pela CLI (criar recurso, deploy, start/stop, etc.) é registrada
localmente em `~/.cmx/activity.log` no formato JSONL, agrupado por dia,
com status (✓/✗), duração e detalhes.

Não é telemetria — nada é enviado para fora da sua máquina.

### Autocomplete para o shell

```sh
cmx completion bash                  # bash
cmx completion zsh                   # zsh
cmx completion fish                  # fish
cmx completion powershell            # PowerShell
```

Para ativar no shell atual:

```sh
source <(cmx completion bash)        # bash
eval "$(cmx completion zsh)"         # zsh
cmx completion fish | source         # fish
```

Após ativar, use **Tab** para autocompletar nomes de apps, bancos e comandos.

## Dicas

- Todos os comandos aceitam **nome ou UUID**. Se o nome for ambíguo, o CMX lista as opções e pede o UUID.
- `NO_COLOR=1` desativa as cores (útil em logs de CI).
O fluxo típico de um deploy: `cmx apps deploy meu-app --wait` → `cmx status`.

Use **Tab** para autocompletar nomes de apps e bancos de dados (requer `cmx completion <shell>`).

## Desenvolvimento

```sh
git clone https://github.com/Vime-Labs/cmx
cd cmx
go build .          # compila para a plataforma atual
make test           # roda os testes
make build-all      # compila para todos os alvos em dist/
make install        # instala em /usr/local/bin (macOS/Linux)
```

Para releases, basta criar uma tag:

```sh
git tag v1.0.0
git push origin v1.0.0
```

O GitHub Actions compila, gera checksums e publica a release automaticamente.
