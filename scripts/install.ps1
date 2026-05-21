# install.ps1 — baixa e instala o cmx da última release do GitHub (Windows)
# Uso: irm https://raw.githubusercontent.com/Vime-Labs/cmx/main/scripts/install.ps1 | iex
# Ou:  .\scripts\install.ps1

$repo       = "Vime-Labs/cmx"
$installDir = "$env:USERPROFILE\.cmx\bin"
$asset      = "cmx-windows-amd64.exe"
$binary     = "$installDir\cmx.exe"

# ── verifica gh CLI ───────────────────────────────────────────────────────────
if (-not (Get-Command gh -ErrorAction SilentlyContinue)) {
    Write-Error @"
GitHub CLI (gh) não encontrado.
Instale em: https://cli.github.com

Ou baixe manualmente e coloque em $installDir`:
  https://github.com/$repo/releases/latest
"@
    exit 1
}

$authCheck = gh auth status 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Error "Você não está autenticado no GitHub CLI. Execute: gh auth login"
    exit 1
}

# ── download ──────────────────────────────────────────────────────────────────
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

Write-Host "Baixando $asset..."
gh release download --repo $repo --pattern $asset --output $binary --clobber

if ($LASTEXITCODE -ne 0) {
    Write-Error "Download falhou."
    exit 1
}

# ── adiciona ao PATH do usuário ───────────────────────────────────────────────
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$userPath;$installDir", "User")
    # atualiza o PATH da sessão atual também
    $env:PATH = "$env:PATH;$installDir"
    Write-Host "Adicionado $installDir ao PATH."
} else {
    Write-Host "$installDir já está no PATH."
}

Write-Host ""
Write-Host "✓ cmx instalado em $binary"
& $binary version
Write-Host ""
Write-Host "Pronto. Execute: cmx --help"
