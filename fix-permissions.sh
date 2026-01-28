#!/bin/bash
# fix-permissions.sh

echo "🔧 Corrigindo permissões do projeto..."

# Remover diretório tmp antigo (se existir)
if [ -d "tmp" ]; then
    echo "🗑️  Removendo diretório tmp antigo..."
    rm -rf tmp
fi

# Criar nova estrutura tmp
echo "📁 Criando nova estrutura de diretórios..."
mkdir -p tmp
chmod 755 tmp

# Corrigir permissões dos scripts
echo "🔐 Corrigindo permissões dos scripts..."
chmod +x dev.sh
chmod +x fix-permissions.sh

# Corrigir permissões do Go
echo "⚙️  Verificando configuração Go..."
if [ -f "go.mod" ]; then
    echo "📦 Atualizando módulos Go..."
    go mod tidy
    go mod download
fi

# Verificar e instalar Air
echo "🔄 Verificando Air..."
if ! command -v air &> /dev/null; then
    echo "⬇️  Instalando Air..."
    go install github.com/air-verse/air@latest
fi

echo "✅ Permissões corrigidas!"
echo ""
echo "Agora execute:"
echo "  ./dev.sh"
echo "ou"
echo "  make dev"
