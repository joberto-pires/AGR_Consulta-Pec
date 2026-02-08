package handlers

import (
	"database/sql"
	"log"
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
    ordenarPor := r.URL.Query().Get("ordenar_por")
    direcao := r.URL.Query().Get("direcao")

    // Valores padrão para ordenação
    if ordenarPor == "" {
        ordenarPor = "id"
    }
    if direcao == "" {
        direcao = "DESC"
    }

    // Validar colunas para evitar SQL injection
    colunasValidas := map[string]bool{
        "id":            true,
        "nome":          true,
        "email":         true,
        "telefone":      true,
        "cpf_cnpj":      true,
        "data_cadastro": true,
    }
    
    if !colunasValidas[ordenarPor] {
        ordenarPor = "id"
    }
    
    // Validar direção
    if direcao != "ASC" && direcao != "DESC" {
        direcao = "DESC"
    }

    // Construir a query
    query := "SELECT id, nome, email, telefone, cpf_cnpj, data_cadastro FROM clientes"
    args := []interface{}{}

    if busca != "" {
        query += " WHERE nome LIKE ? OR cpf_cnpj LIKE ? OR email LIKE ? OR telefone LIKE ?"
        likeBusca := "%" + busca + "%"
        args = append(args, likeBusca, likeBusca, likeBusca, likeBusca)
    }

    query += " ORDER BY " + ordenarPor + " " + direcao + " LIMIT ? OFFSET ?"
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
        countQuery += " WHERE nome LIKE ? OR cpf_cnpj LIKE ? OR email LIKE ? OR telefone LIKE ?"
        countArgs = append(countArgs, "%"+busca+"%", "%"+busca+"%", "%"+busca+"%", "%"+busca+"%")
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

    // Calcular páginas para mostrar
    paginas := calcularPaginacao(pagina, totalPaginas)

    data := map[string]interface{}{
        "Clientes":       clientes,
        "PaginaAtual":    pagina,
        "TotalPaginas":   totalPaginas,
        "TotalRegistros": total,
        "Busca":          busca,
        "OrdenarPor":     ordenarPor,
        "Direcao":        direcao,
        "Paginas":        paginas,
        "Title":          "Clientes",
    }

    // Se for uma requisição HTMX (busca, paginação ou ordenação), renderizar apenas a tabela
    if r.Header.Get("HX-Request") == "true" {
        app.renderTemplate(w, r, "clientes/tabela.html", data)
        return
    }

    // Caso contrário, renderizar a página completa
    app.renderTemplate(w, r, "clientes/lista.html", data)
}

func calcularPaginacao(paginaAtual, totalPaginas int) []int {
    // Mostrar no máximo 5 páginas
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

        row := app.DB.QueryRow("SELECT id, nome, email, telefone, cpf_cnpj FROM clientes WHERE id = ?", id)
        err = row.Scan(&cliente.ID, &cliente.Nome, &cliente.Email, &cliente.Telefone, &cliente.CpfCnpj)
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

  log.Printf("🔍 DEBUG: FormCliente chamado, renderizando clientes/formulario.html")
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

    var result sql.Result
    var err error
    
    if (id == "") || (id=="0") {
        // Inserir novo cliente
        result, err = app.DB.Exec(
            `INSERT INTO clientes 
            (nome, email, telefone, cpf_cnpj, endereco, cidade, estado, observacoes) 
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
            nome, email, telefone, cpfCnpj, endereco, cidade, estado, observacoes,
        )
        if err != nil {
            log.Printf("❌ Erro ao inserir cliente: %v", err)
            app.serverError(w, r, err)
            return
        }
        
        // Obter o ID gerado
        rowsAffected, _ := result.RowsAffected()
        lastID, _ := result.LastInsertId()
        
        log.Printf("✅ Cliente inserido - Rows: %d, LastID: %d", rowsAffected, lastID)
        
        // Verificar se conseguiu obter o ID
        if lastID == 0 {
            // Tentar método alternativo para DuckDB
            var newID int
            err = app.DB.QueryRow("SELECT last_insert_id()").Scan(&newID)
            if err != nil {
                // Outro método
                err = app.DB.QueryRow("SELECT currval('clientes_id_seq')").Scan(&newID)
                if err != nil {
                    log.Printf("⚠️ Não foi possível obter último ID: %v", err)
                } else {
                    log.Printf("✅ ID obtido via currval: %d", newID)
                }
            } else {
                log.Printf("✅ ID obtido via last_insert_id: %d", newID)
            }
        }
    } else {
        // Atualizar cliente existente
        idInt, err := strconv.Atoi(id)
        if err != nil {
            app.serverError(w, r, err)
            return
        }
        
        result, err = app.DB.Exec(
            `UPDATE clientes SET 
            nome=?, email=?, telefone=?, cpf_cnpj=?, 
            endereco=?, cidade=?, estado=?, observacoes=? 
            WHERE id=?`,
            nome, email, telefone, cpfCnpj, endereco, cidade, estado, observacoes, idInt,
        )
        if err != nil {
            log.Printf("❌ Erro ao atualizar cliente: %v", err)
            app.serverError(w, r, err)
            return
        }
        
        rowsAffected, _ := result.RowsAffected()
        log.Printf("✅ Cliente atualizado - ID: %s, Rows: %d", id, rowsAffected)
    }

    // Redirecionar para a lista de clientes
    w.Header().Set("HX-Redirect", "/clientes")
    w.WriteHeader(http.StatusOK)
}// DetalhesCliente exibe os detalhes de um cliente
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

    app.renderTemplate(w, r, "clientes/detalhes_sidebar.html", data)
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

    _, err = app.DB.Exec("DELETE FROM clientes WHERE id = ?;", id)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    w.Header().Set("HX-Trigger", `{"showToast": {"message": "Cliente excluído com sucesso.", "type": "success"}}`)
    w.WriteHeader(http.StatusOK)
}


