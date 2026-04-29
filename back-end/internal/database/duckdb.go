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
	db, err := sql.Open("duckdb", connStr)
	if err != nil {
		return nil, err
	}
 
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
 
		`CREATE SEQUENCE IF NOT EXISTS propriedades_id_seq START 1`,
		`CREATE TABLE IF NOT EXISTS propriedades (
			id INTEGER PRIMARY KEY DEFAULT nextval('propriedades_id_seq'),
			cliente_id INTEGER NOT NULL,
			nome TEXT NOT NULL,
			hectares REAL,
			municipio TEXT,
			estado TEXT,
			coordenadas TEXT,
			caracteristica_solo TEXT,
			infraestrutura TEXT,
			historico_culturas TEXT,
			data_cadastro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (cliente_id) REFERENCES clientes(id)
		)`,
 
		`CREATE SEQUENCE IF NOT EXISTS analises_id_seq START 1`,
		`CREATE TABLE IF NOT EXISTS analises (
			id INTEGER PRIMARY KEY DEFAULT nextval('analises_id_seq'),
			propriedade_id INTEGER NOT NULL,
			talhao TEXT,
			data_amostra DATE,
			ph REAL DEFAULT 0,
			materia_organica REAL DEFAULT 0,
			nitrogenio REAL DEFAULT 0,
			fosforo REAL DEFAULT 0,
			potassio REAL DEFAULT 0,
			ctc REAL DEFAULT 0,
			saturacao_bases REAL DEFAULT 0,
			rec_calcario TEXT,
			rec_adubacao TEXT,
			plano_correcao TEXT,
			status TEXT DEFAULT 'pendente',
			data_cadastro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (propriedade_id) REFERENCES propriedades(id)
		)`,
 
		`CREATE SEQUENCE IF NOT EXISTS consultas_id_seq START 1`,
		`CREATE TABLE IF NOT EXISTS consultas (
			id INTEGER PRIMARY KEY DEFAULT nextval('consultas_id_seq'),
			cliente_id INTEGER NOT NULL,
			propriedade_id INTEGER,
			data_consulta DATE NOT NULL,
			tipo TEXT,
			diagnostico TEXT,
			recomendacoes TEXT,
			planejamento_safra TEXT,
			controle_pragas TEXT,
			gestao_irrigacao TEXT,
			analise_custos TEXT,
			status TEXT DEFAULT 'agendada',
			data_cadastro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (cliente_id) REFERENCES clientes(id),
			FOREIGN KEY (propriedade_id) REFERENCES propriedades(id)
		)`,
 
		`CREATE SEQUENCE IF NOT EXISTS monitoramento_id_seq START 1`,
		`CREATE TABLE IF NOT EXISTS monitoramento (
			id INTEGER PRIMARY KEY DEFAULT nextval('monitoramento_id_seq'),
			propriedade_id INTEGER NOT NULL,
			data_registro DATE NOT NULL,
			tipo_registro TEXT,
			descricao TEXT,
			registro_aplicacao TEXT,
			produtividade REAL DEFAULT 0,
			alertas TEXT,
			data_cadastro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (propriedade_id) REFERENCES propriedades(id)
		)`,
	}
 
	for i, tableSQL := range tables {
		parts := strings.Fields(tableSQL)
		label := "SQL"
		if len(parts) >= 3 {
			label = parts[2]
		}
		log.Printf("📝 Criando: %s", label)
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

