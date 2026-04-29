package main

import (
	"AGR_Consulta-Pec/back-end/internal/database"
	"AGR_Consulta-Pec/back-end/internal/handlers"
	"AGR_Consulta-Pec/back-end/internal/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}
  os.Setenv("DB_PATH", "AGRConsultaPec.db")	
	log.Printf("🚀 Iniciando AgroConsultoria v1.0.0")
	log.Printf("📁 Ambiente: %s", env)
	
	dbPath := "AGRConsultaPec.db"
	if customPath := os.Getenv("DB_PATH"); customPath != "" {
		dbPath = customPath
	}
	
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()
	
	log.Printf("🗄️  Banco de dados conectado: %s", dbPath)
	
  workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ Erro ao obter diretório de trabalho: %v", err)
	}

	log.Printf("📂 Diretório de trabalho atual: %s", workingDir)
	
	templatesPath := filepath.Join(workingDir, "front-end", "templates")
	staticPath := filepath.Join(workingDir, "front-end", "static")
	
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		parentDir := filepath.Dir(workingDir)
		templatesPath = filepath.Join(parentDir, "front-end", "templates")
		staticPath = filepath.Join(parentDir, "front-end", "static")
		
		if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
			grandparentDir := filepath.Dir(parentDir)
			templatesPath = filepath.Join(grandparentDir, "front-end", "templates")
			staticPath = filepath.Join(grandparentDir, "front-end", "static")
		}
	}
	
	log.Printf("🔍 Procurando templates em: %s", templatesPath)
	
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		log.Printf("⚠️  Diretório de templates não encontrado: %s", templatesPath)
	}

	app := &handlers.Application{
		DB:            db.DB,
		TemplatesFS:   templatesPath,
		Env:           env,
		StaticFS:      staticPath,
		StartTime:     time.Now(),
	}
	
	err = app.InitTemplates()
		if err != nil {
			log.Fatalf("❌ Erro ao inicializar templates: %v", err)
		}
		
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

handler := app.Routes()
handler = middleware.CacheMiddleware(handler)
handler = middleware.NoCompressionMiddleware(handler) // Use este em vez de GzipMiddleware

server := &http.Server{
	Addr:         ":" + port,
	Handler:      handler,
	ReadTimeout:  5 * time.Second,   // Reduzido
	WriteTimeout: 10 * time.Second,  // Reduzido
	IdleTimeout:  30 * time.Second,
}
	log.Printf("🌐 Servidor iniciado em http://localhost:%s", port)
	log.Printf("📊 Acesse http://localhost:%s/dashboard para começar", port)
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}
