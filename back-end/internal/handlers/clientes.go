package handlers

import (
	"database/sql"
	"embed"
	"net/http"
	"strconv"
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
  Init(templates, db)
 //Clientes
 http.HandleFunc("/", Homepage)
// DB = db
 http.HandleFunc("/clientes", GetClients) 

 http.HandleFunc("/clientes/New", PostClients) 

 http.HandleFunc("/clientes/Filter", FilterClients)

 http.HandleFunc("/clientes/Edit", PutClients)

 http.HandleFunc("/clientes/Delete", DelClients)
}

func GetClients(w http.ResponseWriter, r *http.Request) {
	 if r.Method == "GET" {
		 rows, err := DB.Query("Select * from clientes order by clientes.nome")
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
		 
		  Tmpl.ExecuteTemplate(w,"clientes/lista.html", clientes)
	 }
}

func PostClients(w http.ResponseWriter, r *http.Request) {
		 if r.Method == "POST" {
			 
			 nome := r.FormValue("nome")
			 email := r.FormValue("email")
			 telefone := r.FormValue("telefone")
			 cpfcnpj := r.FormValue("cpf_cnpj")

			 _, err := DB.Exec("INSERT INTO clientes (nome, email, telefone, cpf_cnpj) VALUES(?,?,?,?)", 
			 nome, email, telefone, cpfcnpj)
			 if err != nil {
				 http.Error(w, err.Error(), http.StatusInternalServerError)
				 return
			 }
			 w.Header().Set("HX-Refresh", "true")
			 w.WriteHeader(http.StatusCreated)
		 }

}

func FilterClients(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/clientes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalido!", http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		var c Cliente
		err :=  DB.QueryRow("Select * from clientes where id = ?", id).
	  Scan(&c.ID, &c.Nome, &c.Email, &c.Telefone, &c.Cpf_cnpj, &c.DataCadastro, &c.UsuarioCadastro)
		if err !=  nil {
	    http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		Tmpl.ExecuteTemplate(w, "clientes/formulario.html", c)
  }
}


func PutClients(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/clientes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.Method == "PUT" {

			 nome := r.FormValue("nome")
			 email := r.FormValue("email")
			 telefone := r.FormValue("telefone")
			 cpfcnpj := r.FormValue("cpf_cnpj")

			 _, err := DB.Exec("UPDATE clientes SET nome=?, email=?, telefone=?, cpf_cnpj=? where id=?", 
			 nome, email, telefone, cpfcnpj, id)
			 if err != nil {
				 http.Error(w, err.Error(), http.StatusInternalServerError)
				 return
			 }
			 w.WriteHeader(http.StatusOK)
	}
}


func DelClients(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/clientes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.Method == "DELETE" {

			 _, err := DB.Exec("DELETE FROM clientes where id=?", id)
			 if err != nil {
				 http.Error(w, err.Error(), http.StatusInternalServerError)
				 return
			 }
			 w.WriteHeader(http.StatusOK)
	}
}













