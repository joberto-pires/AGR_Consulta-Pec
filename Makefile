.PHONY: dev build run clean setup test migrate help

BINARY_NAME = agroconsultoria
BUILD_DIR = bin
TMP_DIR = tmp

# Configuração inicial do projeto
setup:
	@echo "📁 Criando estrutura de diretórios..."
	@mkdir -p front-end/static/{css,js,images}
	@mkdir -p front-end/templates/{clientes,propriedades,consultas,analises,relatorios,components}
	@mkdir -p back-end/{cmd,internal/{handlers,database,models,services},pkg/utils}
	@mkdir -p $(BUILD_DIR) $(TMP_DIR)
	@echo "✅ Estrutura criada"
	
	@if [ ! -f "go.mod" ]; then \
		echo "📦 Inicializando módulo Go..."; \
		go mod init agroconsultoria; \
	fi
	
	@echo "📥 Instalando dependências..."
	@go mod tidy
	@go install github.com/air-verse/air@latest
	@echo "✨ Configuração concluída!"

# Desenvolvimento com live reload
dev:
	@if ! command -v air &> /dev/null; then \
		echo "📦 Instalando Air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@echo "🚀 Iniciando servidor de desenvolvimento..."
	@air -c .air.toml

# Compilar para produção
build:
	@echo "🔨 Compilando aplicação..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./back-end/cmd/main.go
	@echo "✅ Binário criado: $(BUILD_DIR)/$(BINARY_NAME)"

# Executar aplicação compilada
run: build
	@echo "▶️  Executando aplicação..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Limpar arquivos temporários
clean:
	@echo "🧹 Limpando arquivos temporários..."
	@rm -rf $(BUILD_DIR) $(TMP_DIR)
	@find . -name "*.log" -type f -delete
#	@find . -name "*.db" -type f -delete
	@echo "✅ Limpeza concluída"

# Executar testes
test:
	@echo "🧪 Executando testes..."
	@go test ./back-end/... -v

# Criar banco de dados e tabelas
migrate:
	@echo "🗄️  Criando banco de dados..."
	@go run back-end/cmd/migrate.go

# Instalar/atualizar dependências
deps:
	@echo "📦 Atualizando dependências..."
	@go mod tidy
	@go mod download

# Mostrar ajuda
help:
	@echo "Comandos disponíveis:"
	@echo ""
	@echo "  make setup   - Configurar estrutura inicial do projeto"
	@echo "  make dev     - Iniciar servidor com live reload (Air)"
	@echo "  make build   - Compilar aplicação para produção"
	@echo "  make run     - Compilar e executar aplicação"
	@echo "  make clean   - Limpar arquivos temporários"
	@echo "  make test    - Executar testes"
	@echo "  make migrate - Criar banco de dados e tabelas"
	@echo "  make deps    - Atualizar dependências"
	@echo "  make help    - Mostrar esta mensagem"
	@echo ""
	@echo "Para desenvolvimento, use: make dev ou ./dev.sh"

.DEFAULT_GOAL := help
