package handlers

import (
	"database/sql"
	"embed"
	"html/template"
	"net/http"
)


type Cliente struct {
	ID									int						`json:"id"`
	Nome								string		  	`json:"nome"`
	Email								string				`json:"email"`
	Telefone						string				`json:"telefone"`
	Cpf_cnpj						string				`json:"cpf_cnpj"`
	DataCadastro			 	string				`json:"data_cadastro"`
	UsuarioCadastro			int						`json:"usuario_cadastro"`

}


func SetupRoutes(db *sql.DB, templates embed.FS) {
 tmpl := template.Must(template.New("").ParseFS(templates, "frontend/templates/**/*.html"))

 //Clientes
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	 tmpl.ExecuteTemplate(w, "index.html", nil)
 })

 http.HandleFunc("/clientes", func(w http.ResponseWriter, r *http.Request) {
	 if r.Method == "GET" {
		 rows, err := db.Query("Select * from clientes order by clinetes.nome")
		 if err != nil {
			 http.Error(w, err.Error(), http.StatusInternalServerError)
			 return
		 }
		 defer rows.Close()

		 var clientes []Cliente
		 for rows.Next() {
			 var c Cliente
			 err := rows.Scan(&c.ID, &c.Nome, &c.Email, &c.Telefone, &c.Cpf_cnpj, &c.DataCadastro, &c.UsuarioCadastro)
  		 if err != nil {
				 http.Error(w, err.Error(), http.StatusInternalServerError)
				 return
			 }
			 clientes = append(clientes, c)
		 }
		 
		  tmpl.ExecuteTemplate(w,"clientes/lista.html", clientes)

	 } else if r.Method == "POST" {
		 nome := r.FormValue("nome")
		 email := r.FormValue("email")
     telefone := r.FormValue("telefone")
		 cpfcnpj := r.FormValue("cpf_cnpj")

		 _, err := db.Exec("INSERT INTO clientes (nome, email, telefone, cpf_cnpj) VALUES(?,?,?,?)", 
		 nome, email, telefone, cpfcnpj)
		 if err != nil {
			 http.Error(w, err.Error(), http.StatusInternalServerError)
			 return
		 }
	 }
 })
}
