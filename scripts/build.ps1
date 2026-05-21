# build.ps1 — cross-compilation para Windows (sem dependência de make)
# Uso: .\scripts\build.ps1 [-Version v1.0.0]

param(
    [string]$Version = ""
)

if ($Version -eq "") {
    $Version = git describe --tags --always --dirty 2>$null
    if (-not $Version) { $Version = "dev" }
}

$ldflags = "-X github.com/Vime-Labs/cmx/cmd.Version=$Version -s -w"
$dist    = "dist"

New-Item -ItemType Directory -Force -Path $dist | Out-Null

$targets = @(
    @{ GOOS = "linux";   GOARCH = "amd64"; Out = "cmx-linux-amd64"      },
    @{ GOOS = "linux";   GOARCH = "arm64"; Out = "cmx-linux-arm64"      },
    @{ GOOS = "darwin";  GOARCH = "amd64"; Out = "cmx-darwin-amd64"     },
    @{ GOOS = "darwin";  GOARCH = "arm64"; Out = "cmx-darwin-arm64"     },
    @{ GOOS = "windows"; GOARCH = "amd64"; Out = "cmx-windows-amd64.exe"}
)

foreach ($t in $targets) {
    $env:GOOS   = $t.GOOS
    $env:GOARCH = $t.GOARCH
    $out = "$dist/$($t.Out)"
    Write-Host "Building $($t.GOOS)/$($t.GOARCH)..." -NoNewline
    go build -ldflags $ldflags -o $out .
    if ($LASTEXITCODE -eq 0) { Write-Host " OK" } else { Write-Host " FALHOU"; exit 1 }
}

$env:GOOS   = ""
$env:GOARCH = ""

Write-Host "`nBinários em $dist/:"
Get-ChildItem $dist | Format-Table Name, @{L="Tamanho";E={"{0:N0} KB" -f ($_.Length/1KB)}}
