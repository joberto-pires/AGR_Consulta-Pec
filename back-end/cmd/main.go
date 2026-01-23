package main

import (
	"AGR_Consulta-Pec/back-end/internal/database"
//	"embed"
	"log"
)

//var templates embed.FS

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Erro ao Conectar com o Banco de Dados!")
	}
	defer db.Close()

}
