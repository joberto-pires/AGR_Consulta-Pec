package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/marcboeker/go-duckdb"
)

type Database struct {
	*sql.DB
}

func InitDB(dbPath string) (*Database, error) {
	connStr := fmt.Sprintf("%s?access_mode=READ_WRITE&threads=6", dbPath)
 db, err := sql.Open("duckdb",connStr)
 if err != nil {
	 return nil, err
 }

 //teste conexão 
 if err := db.Ping(); err != nil {
	 return nil, fmt.Errorf("erro ao conectar com o banco: %w", err)
 }

 err = createTables(db)
 if err != nil {
	 return nil, err
 }
 return &Database{db}, nil
}


func createTables(db *sql.DB) error {
	tables := []string{
		// Tabela clientes com SERIAL para auto-increment
		`CREATE SEQUENCE IF NOT EXISTS clientes_id_seq START 1`,
		
		`CREATE TABLE IF NOT EXISTS clientes (
			id INTEGER PRIMARY KEY DEFAULT nextval('clientes_id_seq'),
			nome TEXT NOT NULL,
			email TEXT,
			telefone TEXT,
			cpf_cnpj TEXT UNIQUE,
			data_cadastro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			endereco TEXT,
			cidade TEXT,
			estado TEXT,
			observacoes TEXT,
			ativo BOOLEAN DEFAULT true
		)`,
		
		// Tabela propriedades com SERIAL
		`CREATE SEQUENCE IF NOT EXISTS propriedades_id_seq START 1`,
		
		`CREATE TABLE IF NOT EXISTS propriedades (
			id INTEGER PRIMARY KEY DEFAULT nextval('propriedades_id_seq'),
			cliente_id INTEGER NOT NULL,
			nome TEXT NOT NULL,
			hectares REAL,
			municipio TEXT,
			estado TEXT,
			coordenadas TEXT,
			FOREIGN KEY (cliente_id) REFERENCES clientes(id)
		)`,
		
		// Tabela consultas com SERIAL
		`CREATE SEQUENCE IF NOT EXISTS consultas_id_seq START 1`,
		
		`CREATE TABLE IF NOT EXISTS consultas (
			id INTEGER PRIMARY KEY DEFAULT nextval('consultas_id_seq'),
			cliente_id INTEGER NOT NULL,
			propriedade_id INTEGER,
			data_consulta DATE NOT NULL,
			tipo_consulta TEXT,
			observacoes TEXT,
			resultado TEXT,
			FOREIGN KEY (cliente_id) REFERENCES clientes(id),
			FOREIGN KEY (propriedade_id) REFERENCES propriedades(id)
		)`,
		
		// Tabela analises com SERIAL
		`CREATE SEQUENCE IF NOT EXISTS analises_id_seq START 1`,
		
		`CREATE TABLE IF NOT EXISTS analises (
			id INTEGER PRIMARY KEY DEFAULT nextval('analises_id_seq'),
			propriedade_id INTEGER NOT NULL,
			tipo_analise TEXT,
			data_amostra DATE,
			resultado TEXT,
			recomendacoes TEXT,
			FOREIGN KEY (propriedade_id) REFERENCES propriedades(id)
		)`,
	}
	
	for i, tableSQL := range tables {
		log.Printf("📝 Criando: %s", strings.Split(tableSQL, " ")[1])
		_, err := db.Exec(tableSQL)
		if err != nil {
			return fmt.Errorf("erro ao executar SQL %d: %v\nSQL: %s", i+1, err, tableSQL)
		}
	}
	
	log.Println("✅ Tabelas criadas com sucesso")
	return nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}
