# cmx

CLI da Vime Labs para gerenciar aplicações e bancos de dados nos servidores remotos via Coolify.

## Instalação

> **Pré-requisito:** [GitHub CLI](https://cli.github.com) autenticado (`gh auth login`).

**macOS / Linux**
```sh
curl -fsSL https://raw.githubusercontent.com/Vime-Labs/cmx/main/scripts/install.sh | sh
```

**Windows (PowerShell)**
```powershell
irm https://raw.githubusercontent.com/Vime-Labs/cmx/main/scripts/install.ps1 | iex
```

Após a instalação, reinicie o terminal e execute `cmx --help`.

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

### Aplicações

```sh
cmx apps list                        # lista todas as aplicações
cmx apps get <uuid|nome>             # detalhes de uma aplicação
cmx apps deploy <uuid|nome>          # dispara deploy
cmx apps deploy <uuid|nome> --force  # deploy sem cache
cmx apps logs <uuid|nome>            # últimas 100 linhas de log
cmx apps logs <uuid|nome> -n 200     # N linhas de log
cmx apps logs <uuid|nome> --follow   # acompanha logs em tempo real
cmx apps start <uuid|nome>
cmx apps stop <uuid|nome>
cmx apps restart <uuid|nome>
cmx apps create                      # wizard interativo
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
cmx dbs start <uuid|nome>
cmx dbs stop <uuid|nome>
cmx dbs restart <uuid|nome>
cmx dbs create                       # wizard interativo
```

### Deployments

```sh
cmx deploy <uuid|nome>               # deploy por nome ou UUID
cmx deploy --tag production          # deploy de todos com a tag
cmx deploy <uuid|nome> --force

cmx deployments list                 # deployments em andamento
cmx deployments get <uuid>           # detalhes de um deployment
cmx deployments history <uuid|nome>  # histórico da aplicação
cmx deployments cancel <uuid>        # cancela deployment ativo
```

### Outros

```sh
cmx version
cmx configure
```

## Dicas

- Todos os comandos aceitam **nome ou UUID**. Se o nome for ambíguo, o CMX lista as opções e pede o UUID.
- `NO_COLOR=1` desativa as cores (útil em logs de CI).
- O fluxo típico de um deploy: `cmx apps deploy meu-app` → `cmx deployments history meu-app` → `cmx apps logs meu-app --follow`.

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
