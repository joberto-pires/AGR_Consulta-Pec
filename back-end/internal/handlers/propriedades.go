package handlers
 
import (
	"database/sql"
	"net/http"
	"strconv"
)
 
func (app *Application) ListaPropriedades(w http.ResponseWriter, r *http.Request) {
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
		ordenarPor = "p.id"
	}
	if direcao == "" {
		direcao = "DESC"
	}
	colunasValidas := map[string]bool{
		"p.id": true, "p.nome": true, "p.municipio": true, "p.hectares": true, "c.nome": true,
	}
	if !colunasValidas[ordenarPor] {
		ordenarPor = "p.id"
	}
	if direcao != "ASC" && direcao != "DESC" {
		direcao = "DESC"
	}
 
	query := `SELECT p.id, p.nome, p.hectares, p.municipio, p.estado, c.nome as cliente_nome, p.data_cadastro
		FROM propriedades p JOIN clientes c ON c.id = p.cliente_id`
	args := []interface{}{}
	if busca != "" {
		query += " WHERE p.nome LIKE ? OR c.nome LIKE ? OR p.municipio LIKE ?"
		lb := "%" + busca + "%"
		args = append(args, lb, lb, lb)
	}
	query += " ORDER BY " + ordenarPor + " " + direcao + " LIMIT ? OFFSET ?"
	args = append(args, limite, offset)
 
	rows, err := app.DB.Query(query, args...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer rows.Close()
 
	var props []Propriedade
	for rows.Next() {
		var p Propriedade
		if err := rows.Scan(&p.ID, &p.Nome, &p.Hectares, &p.Municipio, &p.Estado, &p.ClienteNome, &p.DataCadastro); err != nil {
			app.serverError(w, r, err)
			return
		}
		props = append(props, p)
	}
 
	countQuery := `SELECT COUNT(*) FROM propriedades p JOIN clientes c ON c.id = p.cliente_id`
	countArgs := []interface{}{}
	if busca != "" {
		countQuery += " WHERE p.nome LIKE ? OR c.nome LIKE ? OR p.municipio LIKE ?"
		lb := "%" + busca + "%"
		countArgs = append(countArgs, lb, lb, lb)
	}
	var total int
	app.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	totalPaginas := total / limite
	if total%limite > 0 {
		totalPaginas++
	}
 
	data := map[string]interface{}{
		"Propriedades":   props,
		"PaginaAtual":    pagina,
		"TotalPaginas":   totalPaginas,
		"TotalRegistros": total,
		"Busca":          busca,
		"OrdenarPor":     ordenarPor,
		"Direcao":        direcao,
		"Paginas":        calcularPaginacao(pagina, totalPaginas),
		"Title":          "Propriedades",
	}
 
	if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("partial") == "1" {
		app.renderTemplate(w, r, "propriedades/tabela.html", data)
		return
	}
	app.renderTemplate(w, r, "propriedades/lista.html", data)
}
 
func (app *Application) FormPropriedade(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var prop Propriedade
	title := "Nova Propriedade"
 
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		row := app.DB.QueryRow(`SELECT id, cliente_id, nome, hectares, municipio, estado, coordenadas, caracteristica_solo, infraestrutura, historico_culturas FROM propriedades WHERE id = ?`, id)
		err = row.Scan(&prop.ID, &prop.ClienteID, &prop.Nome, &prop.Hectares, &prop.Municipio, &prop.Estado, &prop.Coordenadas, &prop.CaracteristicaSolo, &prop.Infraestrutura, &prop.HistoricoCulturas)
		if err != nil && err != sql.ErrNoRows {
			app.serverError(w, r, err)
			return
		}
		title = "Editar Propriedade"
	}
 
	clientes, _ := app.listarClientesSelect()
	estados := estadosBR()
	data := map[string]interface{}{
		"Propriedade": prop,
		"Clientes":    clientes,
		"Estados":     estados,
		"Title":       title,
	}
	app.renderTemplate(w, r, "propriedades/editar_sidebar.html", data)
}
 
func (app *Application) SalvarPropriedade(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, r, err)
		return
	}
	id := r.Form.Get("id")
	clienteID, _ := strconv.Atoi(r.Form.Get("cliente_id"))
	nome := r.Form.Get("nome")
	hectares, _ := strconv.ParseFloat(r.Form.Get("hectares"), 64)
	municipio := r.Form.Get("municipio")
	estado := r.Form.Get("estado")
	coordenadas := r.Form.Get("coordenadas")
	carSolo := r.Form.Get("caracteristica_solo")
	infra := r.Form.Get("infraestrutura")
	hist := r.Form.Get("historico_culturas")
 
	if id == "" || id == "0" {
		_, err := app.DB.Exec(`INSERT INTO propriedades (cliente_id, nome, hectares, municipio, estado, coordenadas, caracteristica_solo, infraestrutura, historico_culturas) VALUES (?,?,?,?,?,?,?,?,?)`,
			clienteID, nome, hectares, municipio, estado, coordenadas, carSolo, infra, hist)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else {
		idInt, _ := strconv.Atoi(id)
		_, err := app.DB.Exec(`UPDATE propriedades SET cliente_id=?, nome=?, hectares=?, municipio=?, estado=?, coordenadas=?, caracteristica_solo=?, infraestrutura=?, historico_culturas=? WHERE id=?`,
			clienteID, nome, hectares, municipio, estado, coordenadas, carSolo, infra, hist, idInt)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	w.Header().Set("HX-Redirect", "/propriedades")
	w.WriteHeader(http.StatusOK)
}
 
func (app *Application) DetalhesPropriedade(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	var p Propriedade
	row := app.DB.QueryRow(`SELECT p.id, p.cliente_id, p.nome, p.hectares, p.municipio, p.estado, p.coordenadas, p.caracteristica_solo, p.infraestrutura, p.historico_culturas, p.data_cadastro, c.nome
		FROM propriedades p JOIN clientes c ON c.id = p.cliente_id WHERE p.id = ?`, id)
	err = row.Scan(&p.ID, &p.ClienteID, &p.Nome, &p.Hectares, &p.Municipio, &p.Estado, &p.Coordenadas, &p.CaracteristicaSolo, &p.Infraestrutura, &p.HistoricoCulturas, &p.DataCadastro, &p.ClienteNome)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
	data := map[string]interface{}{
		"Propriedade": p,
		"Title":       "Detalhes da Propriedade",
	}
	app.renderTemplate(w, r, "propriedades/detalhes_sidebar.html", data)
}
 
func (app *Application) ExcluirPropriedade(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	var count int
	app.DB.QueryRow("SELECT COUNT(*) FROM analises WHERE propriedade_id = ?", id).Scan(&count)
	if count > 0 {
		w.Header().Set("HX-Trigger", `{"showToast": {"message": "Não é possível excluir propriedade com análises vinculadas.", "type": "error"}}`)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	app.DB.Exec("DELETE FROM propriedades WHERE id = ?", id)
	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Propriedade excluída com sucesso.", "type": "success"}}`)
	w.WriteHeader(http.StatusOK)
}
 
func (app *Application) listarClientesSelect() ([]Cliente, error) {
	rows, err := app.DB.Query("SELECT id, nome FROM clientes WHERE ativo = true ORDER BY nome")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clientes []Cliente
	for rows.Next() {
		var c Cliente
		rows.Scan(&c.ID, &c.Nome)
		clientes = append(clientes, c)
	}
	return clientes, nil
}
 
func estadosBR() []string {
	return []string{"AC", "AL", "AP", "AM", "BA", "CE", "DF", "ES", "GO", "MA", "MT", "MS", "MG", "PA", "PB", "PR", "PE", "PI", "RJ", "RN", "RS", "RO", "RR", "SC", "SP", "SE", "TO"}
}
 

