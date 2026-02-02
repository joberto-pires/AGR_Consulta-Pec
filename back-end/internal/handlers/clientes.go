package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
)

// Estruturas para Clientes
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

type ClienteResumo struct {
	ID           int    `json:"id"`
	Nome         string `json:"nome"`
	CpfCnpj      string `json:"cpf_cnpj"`
	Telefone     string `json:"telefone"`
	Propriedades int    `json:"propriedades"`
}

func (app *Application) ListaClientes(w http.ResponseWriter, r *http.Request) {
    pagina, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
    if pagina < 1 {
        pagina = 1
    }
    limite := 10
    offset := (pagina - 1) * limite

    busca := r.URL.Query().Get("busca")

    // Construir a query
    query := "SELECT id, nome, email, telefone, cpf_cnpj, data_cadastro FROM clientes"
    args := []interface{}{}

    if busca != "" {
        query += " WHERE nome LIKE ? OR cpf_cnpj LIKE ? OR email LIKE ?"
        likeBusca := "%" + busca + "%"
        args = append(args, likeBusca, likeBusca, likeBusca)
    }

    query += " ORDER BY id DESC LIMIT ? OFFSET ?"
    args = append(args, limite, offset)

    // Executar a query
    rows, err := app.DB.Query(query, args...)
    if err != nil {
        app.serverError(w, r, err)
        return
    }
    defer rows.Close()

    var clientes []Cliente
    for rows.Next() {
        var c Cliente
        err := rows.Scan(&c.ID, &c.Nome, &c.Email, &c.Telefone, &c.CpfCnpj, &c.DataCadastro)
        if err != nil {
            app.serverError(w, r, err)
            return
        }
        clientes = append(clientes, c)
    }

    // Contar total de registros para paginação
    countQuery := "SELECT COUNT(*) FROM clientes"
    countArgs := []interface{}{}
    if busca != "" {
        countQuery += " WHERE nome LIKE ? OR cpf_cnpj LIKE ? OR email LIKE ?"
        countArgs = append(countArgs, "%"+busca+"%", "%"+busca+"%", "%"+busca+"%")
    }

    var total int
    err = app.DB.QueryRow(countQuery, countArgs...).Scan(&total)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    totalPaginas := total / limite
    if total%limite > 0 {
        totalPaginas++
    }

    data := map[string]interface{}{
        "Clientes":     clientes,
        "PaginaAtual":  pagina,
        "TotalPaginas": totalPaginas,
        "Busca":        busca,
        "Title":        "Clientes",
    }

    // Se for uma requisição HTMX (busca ou paginação), renderizar apenas a tabela
    if r.Header.Get("HX-Request") == "true" {
        app.renderTemplate(w, r, "clientes/tabela.html", data)
        return
    }

    // Caso contrário, renderizar a página completa
    app.renderTemplate(w, r, "clientes/lista.html", data)
}

func (app *Application) FormCliente(w http.ResponseWriter, r *http.Request) {
    // Verificar se é edição
    idStr := r.URL.Query().Get("id")
    var cliente Cliente
    var title string

    if idStr != "" {
        id, err := strconv.Atoi(idStr)
        if err != nil {
            app.serverError(w, r, err)
            return
        }

        row := app.DB.QueryRow("SELECT id, nome, email, telefone, cpf_cnpj, endereco, cidade, estado, observacoes FROM clientes WHERE id = ?", id)
        err = row.Scan(&cliente.ID, &cliente.Nome, &cliente.Email, &cliente.Telefone, &cliente.CpfCnpj, &cliente.Endereco, &cliente.Cidade, &cliente.Estado, &cliente.Observacoes)
        if err != nil && err != sql.ErrNoRows {
            app.serverError(w, r, err)
            return
        }
        title = "Editar Cliente"
    } else {
        title = "Novo Cliente"
    }

    // Lista de estados para o select
    estados := []string{"AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA", "MT", "MS", "MG", "PA", "PB", "PR", "PE", "PI", "RJ", "RN", "RS", "RO", "RR", "SC", "SP", "SE", "TO"}

    data := map[string]interface{}{
        "Cliente": cliente,
        "Estados": estados,
        "Title":   title,
    }

    app.renderTemplate(w, r, "clientes/formulario.html", data)
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
    if id == "" {
        // Inserir novo cliente
        _, err = app.DB.Exec(
            "INSERT INTO clientes (nome, email, telefone, cpf_cnpj, endereco, cidade, estado, observacoes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
            nome, email, telefone, cpfCnpj, endereco, cidade, estado, observacoes,
        )
    } else {
        // Atualizar cliente existente
        _, err = app.DB.Exec(
            "UPDATE clientes SET nome=?, email=?, telefone=?, cpf_cnpj=?, endereco=?, cidade=?, estado=?, observacoes=? WHERE id=?",
            nome, email, telefone, cpfCnpj, endereco, cidade, estado, observacoes, id,
        )
    }

    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Redirecionar para a lista de clientes
    w.Header().Set("HX-Redirect", "/clientes")
    w.WriteHeader(http.StatusOK)
}

// DetalhesCliente exibe os detalhes de um cliente
func (app *Application) DetalhesCliente(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    var cliente Cliente
    row := app.DB.QueryRow("SELECT id, nome, email, telefone, cpf_cnpj, data_cadastro, endereco, cidade, estado, observacoes FROM clientes WHERE id = ?", id)
    err = row.Scan(&cliente.ID, &cliente.Nome, &cliente.Email, &cliente.Telefone, &cliente.CpfCnpj, &cliente.DataCadastro, &cliente.Endereco, &cliente.Cidade, &cliente.Estado, &cliente.Observacoes)
    if err != nil {
        if err == sql.ErrNoRows {
            http.NotFound(w, r)
            return
        }
        app.serverError(w, r, err)
        return
    }

    // Buscar propriedades do cliente
    rows, err := app.DB.Query("SELECT id, nome, hectares, municipio, estado FROM propriedades WHERE cliente_id = ?", id)
    if err != nil {
        app.serverError(w, r, err)
        return
    }
    defer rows.Close()

    type Propriedade struct {
        ID        int
        Nome      string
        Hectares  float64
        Municipio string
        Estado    string
    }
    var propriedades []Propriedade
    for rows.Next() {
        var p Propriedade
        err := rows.Scan(&p.ID, &p.Nome, &p.Hectares, &p.Municipio, &p.Estado)
        if err != nil {
            app.serverError(w, r, err)
            return
        }
        propriedades = append(propriedades, p)
    }

    data := map[string]interface{}{
        "Cliente":       cliente,
        "Propriedades":  propriedades,
        "Title":         "Detalhes do Cliente",
    }

    app.renderTemplate(w, r, "clientes/detalhes.html", data)
}

// ExcluirCliente exclui um cliente
func (app *Application) ExcluirCliente(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    // Verificar se o cliente possui propriedades
    var count int
    row := app.DB.QueryRow("SELECT COUNT(*) FROM propriedades WHERE cliente_id = ?", id)
    row.Scan(&count)
    if count > 0 {
        w.Header().Set("HX-Trigger", `{"showToast": {"message": "Não é possível excluir cliente com propriedades vinculadas.", "type": "error"}}`)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    _, err = app.DB.Exec("DELETE FROM clientes WHERE id = ?", id)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    w.Header().Set("HX-Trigger", `{"showToast": {"message": "Cliente excluído com sucesso.", "type": "success"}}`)
    w.WriteHeader(http.StatusOK)
}


