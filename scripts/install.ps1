# install.ps1 — instala cmx no PATH do usuário (Windows)
# Uso: .\scripts\install.ps1

$installDir = "$env:USERPROFILE\.cmx\bin"
$binary     = "$installDir\cmx.exe"

# 1. Compila
Write-Host "Compilando cmx..."
$version = git describe --tags --always --dirty 2>$null
if (-not $version) { $version = "dev" }

$env:GOOS   = "windows"
$env:GOARCH = "amd64"
go build -ldflags "-X github.com/Vime-Labs/cmx/cmd.Version=$version -s -w" -o $binary .
$env:GOOS = ""; $env:GOARCH = ""

if ($LASTEXITCODE -ne 0) {
    Write-Error "Build falhou."
    exit 1
}

# 2. Garante que o diretório existe
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

# 3. Adiciona ao PATH do usuário se ainda não estiver
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$userPath;$installDir", "User")
    Write-Host "Adicionado $installDir ao PATH do usuário."
    Write-Host "Reinicie o terminal para o PATH atualizar."
} else {
    Write-Host "$installDir já está no PATH."
}

Write-Host ""
Write-Host "cmx $version instalado em $binary"
Write-Host "Execute: cmx --help"
