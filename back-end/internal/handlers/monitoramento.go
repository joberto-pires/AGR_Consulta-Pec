package handlers
 
import (
	"database/sql"
	"net/http"
	"strconv"
)
 
func (app *Application) ListaMonitoramento(w http.ResponseWriter, r *http.Request) {
	pagina, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if pagina < 1 {
		pagina = 1
	}
	limite := 10
	offset := (pagina - 1) * limite
	busca := r.URL.Query().Get("busca")
 
	query := `SELECT m.id, m.propriedade_id, p.nome, c.nome, m.data_registro, m.tipo_registro, m.produtividade, m.alertas, m.data_cadastro
		FROM monitoramento m
		JOIN propriedades p ON p.id = m.propriedade_id
		JOIN clientes c ON c.id = p.cliente_id WHERE 1=1`
	args := []interface{}{}
	if busca != "" {
		query += " AND (p.nome LIKE ? OR c.nome LIKE ? OR m.tipo_registro LIKE ?)"
		lb := "%" + busca + "%"
		args = append(args, lb, lb, lb)
	}
	query += " ORDER BY m.data_registro DESC LIMIT ? OFFSET ?"
	args = append(args, limite, offset)
 
	rows, err := app.DB.Query(query, args...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer rows.Close()
 
	var registros []Monitoramento
	for rows.Next() {
		var m Monitoramento
		if err := rows.Scan(&m.ID, &m.PropriedadeID, &m.PropriedadeNome, &m.ClienteNome, &m.DataRegistro, &m.TipoRegistro, &m.Produtividade, &m.Alertas, &m.DataCadastro); err != nil {
			app.serverError(w, r, err)
			return
		}
		registros = append(registros, m)
	}
 
	countQuery := `SELECT COUNT(*) FROM monitoramento m JOIN propriedades p ON p.id = m.propriedade_id JOIN clientes c ON c.id = p.cliente_id WHERE 1=1`
	countArgs := []interface{}{}
	if busca != "" {
		countQuery += " AND (p.nome LIKE ? OR c.nome LIKE ? OR m.tipo_registro LIKE ?)"
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
		"Registros":      registros,
		"PaginaAtual":    pagina,
		"TotalPaginas":   totalPaginas,
		"TotalRegistros": total,
		"Busca":          busca,
		"Paginas":        calcularPaginacao(pagina, totalPaginas),
		"Title":          "Monitoramento",
	}
 
	if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("partial") == "1" {
		app.renderTemplate(w, r, "monitoramento/tabela.html", data)
		return
	}
	app.renderTemplate(w, r, "monitoramento/lista.html", data)
}
 
func (app *Application) FormMonitoramento(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var m Monitoramento
	title := "Novo Registro"
 
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		row := app.DB.QueryRow(`SELECT id, propriedade_id, data_registro, tipo_registro, descricao, registro_aplicacao, produtividade, alertas FROM monitoramento WHERE id = ?`, id)
		err = row.Scan(&m.ID, &m.PropriedadeID, &m.DataRegistro, &m.TipoRegistro, &m.Descricao, &m.RegistroAplicacao, &m.Produtividade, &m.Alertas)
		if err != nil && err != sql.ErrNoRows {
			app.serverError(w, r, err)
			return
		}
		title = "Editar Registro"
	}
 
	props, _ := app.listarPropriedadesSelect()
	tipos := []string{"Desenvolvimento", "Aplicação", "Produtividade", "Alerta de Praga", "Alerta Climático", "Irrigação", "Manutenção", "Outro"}
	data := map[string]interface{}{
		"Registro":     m,
		"Propriedades": props,
		"Tipos":        tipos,
		"Title":        title,
	}
	app.renderTemplate(w, r, "monitoramento/editar_sidebar.html", data)
}
 
func (app *Application) SalvarMonitoramento(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, r, err)
		return
	}
	id := r.Form.Get("id")
	propID, _ := strconv.Atoi(r.Form.Get("propriedade_id"))
	dataRegistro := r.Form.Get("data_registro")
	tipo := r.Form.Get("tipo_registro")
	descricao := r.Form.Get("descricao")
	regAplicacao := r.Form.Get("registro_aplicacao")
	produtividade, _ := strconv.ParseFloat(r.Form.Get("produtividade"), 64)
	alertas := r.Form.Get("alertas")
 
	if id == "" || id == "0" {
		_, err := app.DB.Exec(`INSERT INTO monitoramento (propriedade_id, data_registro, tipo_registro, descricao, registro_aplicacao, produtividade, alertas) VALUES (?,?,?,?,?,?,?)`,
			propID, dataRegistro, tipo, descricao, regAplicacao, produtividade, alertas)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else {
		idInt, _ := strconv.Atoi(id)
		_, err := app.DB.Exec(`UPDATE monitoramento SET propriedade_id=?, data_registro=?, tipo_registro=?, descricao=?, registro_aplicacao=?, produtividade=?, alertas=? WHERE id=?`,
			propID, dataRegistro, tipo, descricao, regAplicacao, produtividade, alertas, idInt)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	w.Header().Set("HX-Redirect", "/monitoramento")
	w.WriteHeader(http.StatusOK)
}
 
func (app *Application) DetalhesMonitoramento(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	var m Monitoramento
	row := app.DB.QueryRow(`SELECT m.id, m.propriedade_id, p.nome, c.nome, m.data_registro, m.tipo_registro, m.descricao, m.registro_aplicacao, m.produtividade, m.alertas, m.data_cadastro
		FROM monitoramento m JOIN propriedades p ON p.id = m.propriedade_id JOIN clientes c ON c.id = p.cliente_id WHERE m.id = ?`, id)
	err = row.Scan(&m.ID, &m.PropriedadeID, &m.PropriedadeNome, &m.ClienteNome, &m.DataRegistro, &m.TipoRegistro, &m.Descricao, &m.RegistroAplicacao, &m.Produtividade, &m.Alertas, &m.DataCadastro)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
	data := map[string]interface{}{
		"Registro": m,
		"Title":    "Detalhes do Monitoramento",
	}
	app.renderTemplate(w, r, "monitoramento/detalhes_sidebar.html", data)
}
 
func (app *Application) ExcluirMonitoramento(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	app.DB.Exec("DELETE FROM monitoramento WHERE id = ?", id)
	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Registro excluído com sucesso.", "type": "success"}}`)
	w.WriteHeader(http.StatusOK)
}
 

