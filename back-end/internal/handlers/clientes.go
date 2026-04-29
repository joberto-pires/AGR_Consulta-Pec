package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Cliente struct {
	ID           int       `json:"id"`
	Nome         string    `json:"nome"`
	Email        string    `json:"email"`
	Telefone     string    `json:"telefone"`
	CpfCnpj      string    `json:"cpf_cnpj"`
	DataCadastro time.Time `json:"data_cadastro"`
	Endereco     string    `json:"endereco"`
	Cidade       string    `json:"cidade"`
	Estado       string    `json:"estado"`
	Observacoes  string    `json:"observacoes"`
	Ativo        bool      `json:"ativo"`
}

func (app *Application) ListaClientes(w http.ResponseWriter, r *http.Request) {
	pagina, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if pagina < 1 {
		pagina = 1
	}
	limite := 10
	offset := (pagina - 1) * limite
	busca := r.URL.Query().Get("busca")
	ordenarPor := r.URL.Query().Get("ordenar_por")
	direcao := r.URL.Query().Get("direcao")

	if ordenarPor == "" {
		ordenarPor = "id"
	}
	if direcao == "" {
		direcao = "DESC"
	}
	colunasValidas := map[string]bool{
		"id": true, "nome": true, "email": true,
		"telefone": true, "cpf_cnpj": true, "data_cadastro": true,
	}
	if !colunasValidas[ordenarPor] {
		ordenarPor = "id"
	}
	if direcao != "ASC" && direcao != "DESC" {
		direcao = "DESC"
	}

	query := "SELECT id, nome, email, telefone, cpf_cnpj, data_cadastro FROM clientes"
	args := []interface{}{}
	if busca != "" {
		query += " WHERE nome LIKE ? OR cpf_cnpj LIKE ? OR email LIKE ? OR telefone LIKE ?"
		lb := "%" + busca + "%"
		args = append(args, lb, lb, lb, lb)
	}
	query += " ORDER BY " + ordenarPor + " " + direcao + " LIMIT ? OFFSET ?"
	args = append(args, limite, offset)

	rows, err := app.DB.Query(query, args...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer rows.Close()

	var clientes []Cliente
	for rows.Next() {
		var c Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.Email, &c.Telefone, &c.CpfCnpj, &c.DataCadastro); err != nil {
			app.serverError(w, r, err)
			return
		}
		clientes = append(clientes, c)
	}

	countQuery := "SELECT COUNT(*) FROM clientes"
	countArgs := []interface{}{}
	if busca != "" {
		countQuery += " WHERE nome LIKE ? OR cpf_cnpj LIKE ? OR email LIKE ? OR telefone LIKE ?"
		lb := "%" + busca + "%"
		countArgs = append(countArgs, lb, lb, lb, lb)
	}
	var total int
	app.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	totalPaginas := total / limite
	if total%limite > 0 {
		totalPaginas++
	}

	data := map[string]interface{}{
		"Clientes":       clientes,
		"PaginaAtual":    pagina,
		"TotalPaginas":   totalPaginas,
		"TotalRegistros": total,
		"Busca":          busca,
		"OrdenarPor":     ordenarPor,
		"Direcao":        direcao,
		"Paginas":        calcularPaginacao(pagina, totalPaginas),
		"Title":          "Clientes",
	}

	if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("partial") == "1" {
		app.renderTemplate(w, r, "clientes/tabela.html", data)
		return
	}
	app.renderTemplate(w, r, "clientes/lista.html", data)
}

func calcularPaginacao(paginaAtual, totalPaginas int) []int {
	var paginas []int
	inicio := paginaAtual - 2
	if inicio < 1 {
		inicio = 1
	}
	fim := inicio + 4
	if fim > totalPaginas {
		fim = totalPaginas
		inicio = fim - 4
		if inicio < 1 {
			inicio = 1
		}
	}
	for i := inicio; i <= fim; i++ {
		paginas = append(paginas, i)
	}
	return paginas
}

func (app *Application) FormCliente(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var cliente Cliente
	title := "Novo Cliente"

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		row := app.DB.QueryRow(
			"SELECT id, nome, email, telefone, cpf_cnpj, COALESCE(endereco,''), COALESCE(cidade,''), COALESCE(estado,''), COALESCE(observacoes,'') FROM clientes WHERE id = ?", id)
		err = row.Scan(&cliente.ID, &cliente.Nome, &cliente.Email, &cliente.Telefone, &cliente.CpfCnpj,
			&cliente.Endereco, &cliente.Cidade, &cliente.Estado, &cliente.Observacoes)
		if err != nil && err != sql.ErrNoRows {
			app.serverError(w, r, err)
			return
		}
		title = "Editar Cliente"
	}

	data := map[string]interface{}{
		"Cliente": cliente,
		"Estados": estadosBR(),
		"Title":   title,
	}
	app.renderTemplate(w, r, "clientes/editar_sidebar.html", data)
}

func (app *Application) SalvarCliente(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, r, err)
		return
	}

	id := r.Form.Get("id")
	nome := r.Form.Get("nome")
	email := r.Form.Get("email")
	telefone := r.Form.Get("telefone")
	cpfCnpj := r.Form.Get("cpf_cnpj")
	endereco := r.Form.Get("endereco")
	cidade := r.Form.Get("cidade")
	estado := r.Form.Get("estado")
	observacoes := r.Form.Get("observacoes")

	var err error

	if id == "" || id == "0" {
		_, err = app.DB.Exec(
			`INSERT INTO clientes (nome, email, telefone, cpf_cnpj, endereco, cidade, estado, observacoes)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			nome, email, telefone, cpfCnpj, endereco, cidade, estado, observacoes,
		)
		if err != nil {
			log.Printf("❌ Erro ao inserir cliente: %v", err)
			app.serverError(w, r, err)
			return
		}
		log.Printf("✅ Cliente inserido: %s", nome)
	} else {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		_, err = app.DB.Exec(
			`UPDATE clientes SET nome=?, email=?, telefone=?, cpf_cnpj=?,
			 endereco=?, cidade=?, estado=?, observacoes=? WHERE id=?`,
			nome, email, telefone, cpfCnpj, endereco, cidade, estado, observacoes, idInt,
		)
		if err != nil {
			log.Printf("❌ Erro ao atualizar cliente: %v", err)
			app.serverError(w, r, err)
			return
		}
		log.Printf("✅ Cliente atualizado: ID %s", id)
	}

	w.Header().Set("HX-Redirect", "/clientes")
	w.WriteHeader(http.StatusOK)
}

func (app *Application) DetalhesCliente(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var c Cliente
	row := app.DB.QueryRow(
		`SELECT id, nome, email, telefone, cpf_cnpj, data_cadastro,
		 COALESCE(endereco,''), COALESCE(cidade,''), COALESCE(estado,''), COALESCE(observacoes,'')
		 FROM clientes WHERE id = ?`, id)
	err = row.Scan(&c.ID, &c.Nome, &c.Email, &c.Telefone, &c.CpfCnpj, &c.DataCadastro,
		&c.Endereco, &c.Cidade, &c.Estado, &c.Observacoes)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}

	// Propriedades do cliente
	type PropResumo struct {
		ID        int
		Nome      string
		Hectares  float64
		Municipio string
		Estado    string
	}
	rows, err := app.DB.Query(
		"SELECT id, nome, hectares, COALESCE(municipio,''), COALESCE(estado,'') FROM propriedades WHERE cliente_id = ?", id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer rows.Close()
	var props []PropResumo
	for rows.Next() {
		var p PropResumo
		rows.Scan(&p.ID, &p.Nome, &p.Hectares, &p.Municipio, &p.Estado)
		props = append(props, p)
	}

	data := map[string]interface{}{
		"Cliente":      c,
		"Propriedades": props,
		"Title":        "Detalhes do Cliente",
	}
	app.renderTemplate(w, r, "clientes/detalhes_sidebar.html", data)
}

func (app *Application) ExcluirCliente(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var count int
	app.DB.QueryRow("SELECT COUNT(*) FROM propriedades WHERE cliente_id = ?", id).Scan(&count)
	if count > 0 {
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Não é possível excluir cliente com propriedades vinculadas.", "type": "error"}}`)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	app.DB.Exec("DELETE FROM clientes WHERE id = ?", id)
	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Cliente excluído com sucesso.", "type": "success"}}`)
	w.WriteHeader(http.StatusOK)
}
