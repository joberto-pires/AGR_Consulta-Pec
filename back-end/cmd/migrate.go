package main

import (
	"AGR_Consulta-Pec/back-end/internal/database"
	"log"
)

func migrate() {
	log.Println("🗄️  Executando migrações do banco de dados...")
	
	db, err := database.InitDB("agroconsultoria.db")
	if err != nil {
		log.Fatalf("❌ Erro: %v", err)
	}
	defer db.Close()
	
	log.Println("✅ Banco de dados criado e migrado com sucesso!")
	log.Println("📊 Tabelas disponíveis:")
	
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("❌ Erro ao listar tabelas: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Printf("⚠️  Erro ao ler tabela: %v", err)
			continue
		}
		log.Printf("   - %s", tableName)
	}

	log.Println("✨ Migração concluída!")
}
