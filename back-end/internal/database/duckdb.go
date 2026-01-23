package database

import (
	"database/sql"
	"fmt"
	_ "github.com/marcboeker/go-duckdb"
)

func InitDB() (*sql.DB, error) {
 db, err := sql.Open("duckdb","AGRConsultaPec.db")
 if err != nil {
	 return nil, err
 }

 err = createTables(db)
 if err != nil {
	 return nil, err
 }
 return db, nil
}


func createTables(db *sql.DB) error {
	queries := []string{
// Clientes
        `CREATE TABLE IF NOT EXISTS clientes (
            id INTEGER PRIMARY KEY,
            nome VARCHAR(100) NOT NULL,
            email VARCHAR(100),
            telefone VARCHAR(20),
            cpf_cnpj VARCHAR(20),
            data_cadastro TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,

        // Propriedades
        `CREATE TABLE IF NOT EXISTS propriedades (
            id INTEGER PRIMARY KEY,
            cliente_id INTEGER,
            nome VARCHAR(100),
            hectares DECIMAL(10,2),
            municipio VARCHAR(100),
            estado VARCHAR(2),
            coordenadas VARCHAR(100),
            tipo_solo VARCHAR(50),
            FOREIGN KEY (cliente_id) REFERENCES clientes(id)
        )`,

        // Talhões
        `CREATE TABLE IF NOT EXISTS talhoes (
            id INTEGER PRIMARY KEY,
            propriedade_id INTEGER,
            numero VARCHAR(10),
            area_hectares DECIMAL(10,2),
            cultura_atual VARCHAR(50),
            data_plantio DATE,
            produtividade_esperada DECIMAL(10,2),
            FOREIGN KEY (propriedade_id) REFERENCES propriedades(id)
        )`,

        // Consultas/Visitas
        `CREATE TABLE IF NOT EXISTS consultas (
            id INTEGER PRIMARY KEY,
            cliente_id INTEGER,
            propriedade_id INTEGER,
            data_consulta DATE,
            tipo_consulta VARCHAR(50),
            observacoes TEXT,
            recomendacoes TEXT,
            proxima_visita DATE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (cliente_id) REFERENCES clientes(id),
            FOREIGN KEY (propriedade_id) REFERENCES propriedades(id)
        )`,

        // Análises de Solo
        `CREATE TABLE IF NOT EXISTS analises_solo (
            id INTEGER PRIMARY KEY,
            talhao_id INTEGER,
            data_analise DATE,
            ph DECIMAL(3,1),
            materia_organica DECIMAL(4,2),
            fosforo_mg_dm3 DECIMAL(6,2),
            potassio_cmolc_dm3 DECIMAL(6,2),
            calcio_cmolc_dm3 DECIMAL(6,2),
            magnesio_cmolc_dm3 DECIMAL(6,2),
            recomendacao_calcario DECIMAL(8,2),
            recomendacao_adubacao TEXT,
            FOREIGN KEY (talhao_id) REFERENCES talhoes(id)
        )`,

        // Registros de Produção
        `CREATE TABLE IF NOT EXISTS producoes (
            id INTEGER PRIMARY KEY,
            talhao_id INTEGER,
            safra VARCHAR(10),
            cultura VARCHAR(50),
            produtividade_ton_ha DECIMAL(8,2),
            data_colheita DATE,
            custo_hectare DECIMAL(10,2),
            receita_hectare DECIMAL(10,2),
            lucro_hectare DECIMAL(10,2)
        )`,		
	}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("Erro ao criar Tabela: %v", err)
		}
	}
   return nil
}
