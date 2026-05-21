#!/bin/sh
# install.sh — baixa e instala o cmx da última release do GitHub
# Uso: curl -fsSL https://raw.githubusercontent.com/Vime-Labs/cmx/main/scripts/install.sh | sh
# Ou:  ./scripts/install.sh
set -e

REPO="Vime-Labs/cmx"
BIN_DIR="/usr/local/bin"

# ── detecta OS ───────────────────────────────────────────────────────────────
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux|darwin) ;;
  *) echo "Sistema operacional não suportado: $OS"; exit 1 ;;
esac

# ── detecta arquitetura ───────────────────────────────────────────────────────
ARCH=$(uname -m)
case "$ARCH" in
  x86_64)          ARCH="amd64" ;;
  aarch64|arm64)   ARCH="arm64" ;;
  *) echo "Arquitetura não suportada: $ARCH"; exit 1 ;;
esac

ASSET="cmx-${OS}-${ARCH}"
echo "Detectado: ${OS}/${ARCH} → ${ASSET}"

# ── verifica gh CLI ───────────────────────────────────────────────────────────
if ! command -v gh >/dev/null 2>&1; then
  echo ""
  echo "GitHub CLI (gh) não encontrado."
  echo "Instale em: https://cli.github.com"
  echo ""
  echo "Ou baixe manualmente e mova para $BIN_DIR:"
  echo "  https://github.com/${REPO}/releases/latest"
  exit 1
fi

if ! gh auth status >/dev/null 2>&1; then
  echo "Você não está autenticado no GitHub CLI. Execute: gh auth login"
  exit 1
fi

# ── download ──────────────────────────────────────────────────────────────────
TMP=$(mktemp)
echo "Baixando ${ASSET}..."
gh release download --repo "$REPO" --pattern "$ASSET" --output "$TMP" --clobber

chmod +x "$TMP"

# ── instala ───────────────────────────────────────────────────────────────────
if [ -w "$BIN_DIR" ]; then
  mv "$TMP" "$BIN_DIR/cmx"
else
  echo "Necessário sudo para instalar em $BIN_DIR"
  sudo mv "$TMP" "$BIN_DIR/cmx"
fi

echo ""
echo "✓ cmx instalado em $BIN_DIR/cmx"
cmx version
