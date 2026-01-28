.PHONY: dev build run clean setup test fix-perms help

BINARY_NAME = agroconsultoria
BUILD_DIR = bin

# Corrigir permissões
fix-perms:
	@echo "🔧 Corrigindo permissões..."
	@chmod +x dev.sh fix-permissions.sh 2>/dev/null || true
	@rm -rf tmp 2>/dev/null || true
	@mkdir -p tmp
	@chmod 755 tmp
	@echo "✅ Permissões corrigidas"

# Setup inicial
setup: fix-perms
	@echo "📁 Criando estrutura..."
	@mkdir -p front-end/static/{css,js,images}
	@mkdir -p front-end/templates/{clientes,propriedades,consultas,analises,relatorios,components}
	@mkdir -p back-end/{cmd,internal/{handlers,database,models,services},pkg/utils}
	@mkdir -p $(BUILD_DIR)
	
	@if [ ! -f "go.mod" ]; then \
		echo "📦 Inicializando módulo Go..."; \
		go mod init agroconsultoria; \
	fi
	
	@echo "📥 Instalando dependências..."
	@go mod tidy
	@go install github.com/air-verse/air@latest
	@echo "✨ Setup concluído!"

# Desenvolvimento
dev: fix-perms
	@echo "🚀 Iniciando desenvolvimento..."
	@if ! command -v air &> /dev/null; then \
		echo "📦 Instalando Air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@air -c .air.toml

# Build
build:
	@echo "🔨 Compilando..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./back-end/cmd/main.go
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)
	@echo "✅ Binário: $(BUILD_DIR)/$(BINARY_NAME)"

# Executar
run: build
	@echo "▶️  Executando..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Limpar
clean:
	@echo "🧹 Limpando..."
	@rm -rf tmp $(BUILD_DIR) 2>/dev/null || true
	@find . -name "*.db" -type f -delete 2>/dev/null || true
	@echo "✅ Limpeza concluída"

# Migrate
migrate:
	@echo "🗄️  Criando banco..."
	@go run back-end/cmd/migrate.go

# Testes
test:
	@echo "🧪 Testando..."
	@go test ./back-end/... -v

# Ajuda
help:
	@echo "📋 Comandos disponíveis:"
	@echo "  make fix-perms  - Corrigir permissões (execute primeiro!)"
	@echo "  make setup      - Configurar projeto completo"
	@echo "  make dev        - Desenvolvimento com live reload"
	@echo "  make build      - Compilar para produção"
	@echo "  make run        - Executar aplicação compilada"
	@echo "  make clean      - Limpar arquivos temporários"
	@echo "  make migrate    - Criar banco de dados"
	@echo "  make test       - Executar testes"
	@echo ""
	@echo "🔧 Solução de problemas:"
	@echo "  Se tiver erro de permissão, execute: make fix-perms"
	@echo "  Se Air não funcionar, execute: make setup"
	@echo ""
	@echo "💡 Dica: Para começar: make setup && make dev"

.DEFAULT_GOAL := help
