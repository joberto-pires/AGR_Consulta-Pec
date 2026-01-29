package main

import (
	FS "AGR_Consulta-Pec"
	"AGR_Consulta-Pec/back-end/internal/database"
	"AGR_Consulta-Pec/back-end/internal/handlers"
	"AGR_Consulta-Pec/back-end/pkg/utils"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
// Configuração do ambiente
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}
	
	log.Printf("🚀 Iniciando AgroConsultoria v1.0.0")
	log.Printf("📁 Ambiente: %s", env)
	
	// Inicializar banco de dados
	dbPath := "AGR_Consulta-Pec.db"
	if customPath := os.Getenv("DB_PATH"); customPath != "" {
		dbPath = customPath
	}
	
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()
	
	log.Printf("🗄️  Banco de dados conectado: %s", dbPath)
	
	// Carregar templates
	tmpl, err := utils.LoadTemplates(FS.TemplatesFS)
	if err != nil {
		log.Fatalf("❌ Erro ao carregar templates: %v", err)
	}
	
	log.Printf("📄 Templates carregados: %d", len(tmpl.Templates()))
	
	// Configurar handlers
	app := &handlers.Application{
		DB:         db.DB,
		Template:   tmpl,
		Env:        env,
		StaticFS:   FS.StaticFS,
		StartTime:  time.Now(),
	}
	
	// Configurar servidor HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      app.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Printf("🌐 Servidor iniciado em http://localhost:%s", port)
	log.Printf("📊 Acesse http://localhost:%s/dashboard para começar", port)
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}
