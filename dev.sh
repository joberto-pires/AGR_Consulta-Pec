#!/bin/bash
# dev.sh - Corrigido

set -e  # Para em caso de erro

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Agro Consultoria - Ambiente de Desenvolvimento${NC}"
echo "======================================================"

# Função para limpar tmp
clean_tmp() {
    echo -e "${BLUE}🧹 Limpando diretório tmp...${NC}"
    rm -rf tmp 2>/dev/null || true
    mkdir -p tmp
    chmod 755 tmp
}

# Limpar tmp no início
clean_tmp

# Verificar Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go não está instalado${NC}"
    exit 1
fi

echo -e "${BLUE}✓ Go $(go version)${NC}"

# Criar estrutura de diretórios
echo -e "${BLUE}📁 Criando estrutura...${NC}"
mkdir -p front-end/static/{css,js,images}
mkdir -p front-end/templates/{clientes,propriedades,consultas,analises,relatorios,components}
mkdir -p back-end/{cmd,internal/{handlers,database,models,services},pkg/utils}
mkdir -p bin

# Configurar módulo Go
if [ ! -f "go.mod" ]; then
    echo -e "${BLUE}📦 Inicializando módulo Go...${NC}"
    go mod init agroconsultoria
fi

# Atualizar dependências
echo -e "${BLUE}📥 Atualizando dependências...${NC}"
go mod tidy
go mod download

# Instalar/verificar Air
if ! command -v air &> /dev/null; then
    echo -e "${BLUE}⬇️  Instalando Air...${NC}"
    go install github.com/air-verse/air@latest
    export PATH="$PATH:$(go env GOPATH)/bin"
fi

echo -e "${BLUE}✓ Air $(air -v 2>/dev/null || echo 'instalado')${NC}"

# Verificar arquivos críticos
echo -e "${BLUE}🔍 Verificando arquivos críticos...${NC}"

# Criar .air.toml se não existir
if [ ! -f ".air.toml" ]; then
    echo -e "${BLUE}📄 Criando .air.toml...${NC}"
    cat > .air.toml << 'EOF'
root = "."
tmp_dir = "tmp"

[build]
  entrypoint = ["back-end/cmd/main.go"]
  cmd = "go build -o ./tmp/main ./back-end/cmd/main.go"
  bin = "./tmp/main"
  delay = 1000
  exclude_dir = ["tmp", ".git", "vendor", "node_modules"]
  include_dir = ["backend", "frontend"]
  include_ext = ["go", "html", "css", "js"]
  send_interrupt = true

[log]
  time = true

[misc]
  clean_on_exit = true
EOF
fi

# Criar main.go básico se não existir
if [ ! -f "back-end/cmd/main.go" ]; then
    echo -e "${BLUE}📄 Criando main.go básico...${NC}"
    mkdir -p back-end/cmd
    cat > back-end/cmd/main.go << 'EOF'
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "🚀 Agro Consultoria está funcionando!")
    })
    
    fmt.Println("Servidor rodando na porta 8080...")
    http.ListenAndServe(":8080", nil)
}
EOF
fi

# Criar CSS básico
if [ ! -f "front-end/static/css/style.css" ]; then
    echo -e "${BLUE}🎨 Criando CSS básico...${NC}"
    cat > front-end/static/css/style.css << 'EOF'
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: system-ui, -apple-system, sans-serif;
    background: #f5f5f5;
    padding: 20px;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    background: white;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

h1 {
    color: #2e7d32;
    margin-bottom: 20px;
}
EOF
fi

# Banco de dados
if [ ! -f "agroconsultoria.db" ]; then
    echo -e "${BLUE}🗄️  Banco de dados não encontrado${NC}"
    echo -e "${BLUE}   Será criado automaticamente${NC}"
fi

echo "======================================================"
echo -e "${GREEN}✅ Ambiente configurado!${NC}"
echo ""
echo -e "${BLUE}Iniciando servidor com live reload...${NC}"
echo -e "${BLUE}Acesse:${NC} http://localhost:8080"
echo "======================================================"

# Executar Air
exec air -c .air.toml
