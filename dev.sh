#!/bin/bash
# dev.sh

GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Agro Consultoria - Ambiente de Desenvolvimento${NC}"
echo "======================================================"

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go não está instalado. Por favor, instale Go 1.21+"
    exit 1
fi

echo -e "${BLUE}✓ Go $(go version)${NC}"

# Verificar estrutura de diretórios
echo -e "${BLUE}📁 Verificando estrutura de diretórios...${NC}"

mkdir -p front-end/static/{css,js,images}
mkdir -p front-end/templates/{clientes,propriedades,consultas,analises,relatorios,components}
mkdir -p back-end/{cmd,internal/{handlers,database,models,services},pkg/utils}
mkdir -p tmp bin

# Verificar dependências
echo -e "${BLUE}📦 Verificando dependências...${NC}"

if [ ! -f "go.mod" ]; then
    go mod init agroconsultoria
fi

go mod tidy

# Instalar Air se necessário
if ! command -v air &> /dev/null; then
    echo -e "${BLUE}⬇️  Instalando Air (live reload)...${NC}"
    go install github.com/air-verse/air@latest
    
    # Adicionar ao PATH se necessário
    if [[ ":$PATH:" != *":$(go env GOPATH)/bin:"* ]]; then
        export PATH="$PATH:$(go env GOPATH)/bin"
        echo "export PATH=\"\$PATH:$(go env GOPATH)/bin\"" >> ~/.bashrc
        echo "export PATH=\"\$PATH:$(go env GOPATH)/bin\"" >> ~/.zshrc
    fi
fi

echo -e "${BLUE}✓ Air $(air -v 2>/dev/null || echo 'instalado')${NC}"

# Criar arquivos iniciais se não existirem
if [ ! -f "front-end/static/css/style.css" ]; then
    echo -e "${BLUE}📄 Criando arquivos CSS...${NC}"
    cat > front-end/static/css/style.css << 'EOF'
/* CSS será criado automaticamente */
EOF
fi

if [ ! -f "front-end/templates/base.html" ]; then
    echo -e "${BLUE}📄 Criando template base...${NC}"
    cat > front-end/templates/base.html << 'EOF'
<!-- Template base será criado -->
EOF
fi

# Verificar banco de dados
if [ ! -f "agroconsultoria.db" ]; then
    echo -e "${BLUE}🗄️  Banco de dados não encontrado${NC}"
    echo -e "${BLUE}   Será criado automaticamente ao iniciar o servidor${NC}"
fi

echo "======================================================"
echo -e "${GREEN}✅ Ambiente configurado com sucesso!${NC}"
echo ""
echo -e "${BLUE}Comandos disponíveis:${NC}"
echo -e "  ${GREEN}./dev.sh${NC}          - Iniciar servidor com live reload"
echo -e "  ${GREEN}make build${NC}        - Compilar para produção"
echo -e "  ${GREEN}make clean${NC}        - Limpar arquivos temporários"
echo ""
echo -e "${BLUE}Acesse:${NC} http://localhost:8080"
echo "======================================================"

# Iniciar Air
air -c .air.toml
