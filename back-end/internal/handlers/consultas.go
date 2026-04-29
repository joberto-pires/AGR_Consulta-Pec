package handlers
 
import (
	"database/sql"
	"net/http"
	"strconv"
)
 
func (app *Application) ListaConsultas(w http.ResponseWriter, r *http.Request) {
	pagina, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if pagina < 1 {
		pagina = 1
	}
	limite := 10
	offset := (pagina - 1) * limite
	busca := r.URL.Query().Get("busca")
	status := r.URL.Query().Get("status")
 
	query := `SELECT co.id, co.cliente_id, cl.nome, COALESCE(p.nome, ''), co.data_consulta, co.tipo, co.status, co.data_cadastro
		FROM consultas co
		JOIN clientes cl ON cl.id = co.cliente_id
		LEFT JOIN propriedades p ON p.id = co.propriedade_id
		WHERE 1=1`
	args := []interface{}{}
	if busca != "" {
		query += " AND (cl.nome LIKE ? OR co.tipo LIKE ?)"
		lb := "%" + busca + "%"
		args = append(args, lb, lb)
	}
	if status != "" {
		query += " AND co.status = ?"
		args = append(args, status)
	}
	query += " ORDER BY co.data_consulta DESC LIMIT ? OFFSET ?"
	args = append(args, limite, offset)
 
	rows, err := app.DB.Query(query, args...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer rows.Close()
 
	var consultas []Consulta
	for rows.Next() {
		var c Consulta
		if err := rows.Scan(&c.ID, &c.ClienteID, &c.ClienteNome, &c.PropriedadeNome, &c.DataConsulta, &c.Tipo, &c.Status, &c.DataCadastro); err != nil {
			app.serverError(w, r, err)
			return
		}
		consultas = append(consultas, c)
	}
 
	countQuery := `SELECT COUNT(*) FROM consultas co JOIN clientes cl ON cl.id = co.cliente_id WHERE 1=1`
	countArgs := []interface{}{}
	if busca != "" {
		countQuery += " AND (cl.nome LIKE ? OR co.tipo LIKE ?)"
		lb := "%" + busca + "%"
		countArgs = append(countArgs, lb, lb)
	}
	if status != "" {
		countQuery += " AND co.status = ?"
		countArgs = append(countArgs, status)
	}
	var total int
	app.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	totalPaginas := total / limite
	if total%limite > 0 {
		totalPaginas++
	}
 
	data := map[string]interface{}{
		"Consultas":      consultas,
		"PaginaAtual":    pagina,
		"TotalPaginas":   totalPaginas,
		"TotalRegistros": total,
		"Busca":          busca,
		"StatusFiltro":   status,
		"Paginas":        calcularPaginacao(pagina, totalPaginas),
		"Title":          "Consultas Técnicas",
	}
 
	if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("partial") == "1" {
		app.renderTemplate(w, r, "consultas/tabela.html", data)
		return
	}
	app.renderTemplate(w, r, "consultas/lista.html", data)
}
 
func (app *Application) FormConsulta(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var c Consulta
	title := "Nova Consulta"
 
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		row := app.DB.QueryRow(`SELECT id, cliente_id, COALESCE(propriedade_id, 0), data_consulta, tipo, diagnostico, recomendacoes, planejamento_safra, controle_pragas, gestao_irrigacao, analise_custos, status FROM consultas WHERE id = ?`, id)
		err = row.Scan(&c.ID, &c.ClienteID, &c.PropriedadeID, &c.DataConsulta, &c.Tipo, &c.Diagnostico, &c.Recomendacoes, &c.PlanejamentoSafra, &c.ControlePragas, &c.GestaoIrrigacao, &c.AnaliseCustos, &c.Status)
		if err != nil && err != sql.ErrNoRows {
			app.serverError(w, r, err)
			return
		}
		title = "Editar Consulta"
	}
 
	clientes, _ := app.listarClientesSelect()
	props, _ := app.listarPropriedadesSelect()
	tiposConsulta := []string{"Visita Técnica", "Análise de Solo", "Diagnóstico", "Planejamento de Safra", "Controle de Pragas", "Reunião com Cliente", "Outro"}
	statusOpts := []string{"agendada", "em_andamento", "concluida", "cancelada"}
 
	data := map[string]interface{}{
		"Consulta":      c,
		"Clientes":      clientes,
		"Propriedades":  props,
		"Tipos":         tiposConsulta,
		"StatusOpts":    statusOpts,
		"Title":         title,
	}
	app.renderTemplate(w, r, "consultas/editar_sidebar.html", data)
}
 
func (app *Application) SalvarConsulta(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, r, err)
		return
	}
	id := r.Form.Get("id")
	clienteID, _ := strconv.Atoi(r.Form.Get("cliente_id"))
	propIDStr := r.Form.Get("propriedade_id")
	dataConsulta := r.Form.Get("data_consulta")
	tipo := r.Form.Get("tipo")
	diagnostico := r.Form.Get("diagnostico")
	recomendacoes := r.Form.Get("recomendacoes")
	planSafra := r.Form.Get("planejamento_safra")
	pragas := r.Form.Get("controle_pragas")
	irrigacao := r.Form.Get("gestao_irrigacao")
	custos := r.Form.Get("analise_custos")
	status := r.Form.Get("status")
	if status == "" {
		status = "agendada"
	}
 
	var propID interface{}
	if propIDStr != "" && propIDStr != "0" {
		propID, _ = strconv.Atoi(propIDStr)
	}
 
	if id == "" || id == "0" {
		_, err := app.DB.Exec(`INSERT INTO consultas (cliente_id, propriedade_id, data_consulta, tipo, diagnostico, recomendacoes, planejamento_safra, controle_pragas, gestao_irrigacao, analise_custos, status) VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
			clienteID, propID, dataConsulta, tipo, diagnostico, recomendacoes, planSafra, pragas, irrigacao, custos, status)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else {
		idInt, _ := strconv.Atoi(id)
		_, err := app.DB.Exec(`UPDATE consultas SET cliente_id=?, propriedade_id=?, data_consulta=?, tipo=?, diagnostico=?, recomendacoes=?, planejamento_safra=?, controle_pragas=?, gestao_irrigacao=?, analise_custos=?, status=? WHERE id=?`,
			clienteID, propID, dataConsulta, tipo, diagnostico, recomendacoes, planSafra, pragas, irrigacao, custos, status, idInt)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	w.Header().Set("HX-Redirect", "/consultas")
	w.WriteHeader(http.StatusOK)
}
 
func (app *Application) DetalhesConsulta(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	var c Consulta
	row := app.DB.QueryRow(`SELECT co.id, co.cliente_id, cl.nome, COALESCE(co.propriedade_id,0), COALESCE(p.nome,''), co.data_consulta, co.tipo, co.diagnostico, co.recomendacoes, co.planejamento_safra, co.controle_pragas, co.gestao_irrigacao, co.analise_custos, co.status, co.data_cadastro
		FROM consultas co JOIN clientes cl ON cl.id = co.cliente_id LEFT JOIN propriedades p ON p.id = co.propriedade_id WHERE co.id = ?`, id)
	err = row.Scan(&c.ID, &c.ClienteID, &c.ClienteNome, &c.PropriedadeID, &c.PropriedadeNome, &c.DataConsulta, &c.Tipo, &c.Diagnostico, &c.Recomendacoes, &c.PlanejamentoSafra, &c.ControlePragas, &c.GestaoIrrigacao, &c.AnaliseCustos, &c.Status, &c.DataCadastro)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
	data := map[string]interface{}{
		"Consulta": c,
		"Title":    "Detalhes da Consulta",
	}
	app.renderTemplate(w, r, "consultas/detalhes_sidebar.html", data)
}
 
func (app *Application) ExcluirConsulta(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	app.DB.Exec("DELETE FROM consultas WHERE id = ?", id)
	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Consulta excluída com sucesso.", "type": "success"}}`)
	w.WriteHeader(http.StatusOK)
}

