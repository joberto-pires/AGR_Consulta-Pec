package handlers
 
import (
	"database/sql"
	"net/http"
	"strconv"
)
 
func (app *Application) ListaAnalises(w http.ResponseWriter, r *http.Request) {
	pagina, _ := strconv.Atoi(r.URL.Query().Get("pagina"))
	if pagina < 1 {
		pagina = 1
	}
	limite := 10
	offset := (pagina - 1) * limite
	busca := r.URL.Query().Get("busca")
	status := r.URL.Query().Get("status")
 
	query := `SELECT a.id, a.propriedade_id, p.nome, c.nome, a.talhao, a.data_amostra, a.ph, a.status, a.data_cadastro
		FROM analises a
		JOIN propriedades p ON p.id = a.propriedade_id
		JOIN clientes c ON c.id = p.cliente_id WHERE 1=1`
	args := []interface{}{}
	if busca != "" {
		query += " AND (p.nome LIKE ? OR c.nome LIKE ? OR a.talhao LIKE ?)"
		lb := "%" + busca + "%"
		args = append(args, lb, lb, lb)
	}
	if status != "" {
		query += " AND a.status = ?"
		args = append(args, status)
	}
	query += " ORDER BY a.id DESC LIMIT ? OFFSET ?"
	args = append(args, limite, offset)
 
	rows, err := app.DB.Query(query, args...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer rows.Close()
 
	var analises []AnaliseSolo
	for rows.Next() {
		var a AnaliseSolo
		if err := rows.Scan(&a.ID, &a.PropriedadeID, &a.PropriedadeNome, &a.ClienteNome, &a.Talhao, &a.DataAmostra, &a.PH, &a.Status, &a.DataCadastro); err != nil {
			app.serverError(w, r, err)
			return
		}
		analises = append(analises, a)
	}
 
	countQuery := `SELECT COUNT(*) FROM analises a JOIN propriedades p ON p.id = a.propriedade_id JOIN clientes c ON c.id = p.cliente_id WHERE 1=1`
	countArgs := []interface{}{}
	if busca != "" {
		countQuery += " AND (p.nome LIKE ? OR c.nome LIKE ? OR a.talhao LIKE ?)"
		lb := "%" + busca + "%"
		countArgs = append(countArgs, lb, lb, lb)
	}
	if status != "" {
		countQuery += " AND a.status = ?"
		countArgs = append(countArgs, status)
	}
	var total int
	app.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	totalPaginas := total / limite
	if total%limite > 0 {
		totalPaginas++
	}
 
	data := map[string]interface{}{
		"Analises":       analises,
		"PaginaAtual":    pagina,
		"TotalPaginas":   totalPaginas,
		"TotalRegistros": total,
		"Busca":          busca,
		"StatusFiltro":   status,
		"Paginas":        calcularPaginacao(pagina, totalPaginas),
		"Title":          "Análises de Solo",
	}
 
	if r.Header.Get("HX-Request") == "true" && r.URL.Query().Get("partial") == "1" {
		app.renderTemplate(w, r, "analises/tabela.html", data)
		return
	}
	app.renderTemplate(w, r, "analises/lista.html", data)
}
 
func (app *Application) FormAnalise(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var a AnaliseSolo
	title := "Nova Análise de Solo"
 
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		row := app.DB.QueryRow(`SELECT id, propriedade_id, talhao, data_amostra, ph, materia_organica, nitrogenio, fosforo, potassio, ctc, saturacao_bases, rec_calcario, rec_adubacao, plano_correcao, status FROM analises WHERE id = ?`, id)
		err = row.Scan(&a.ID, &a.PropriedadeID, &a.Talhao, &a.DataAmostra, &a.PH, &a.MateriaOrganica, &a.Nitrogenio, &a.Fosforo, &a.Potassio, &a.CTC, &a.SaturacaoBases, &a.RecCalcario, &a.RecAdubacao, &a.PlanoCorrecao, &a.Status)
		if err != nil && err != sql.ErrNoRows {
			app.serverError(w, r, err)
			return
		}
		title = "Editar Análise"
	}
 
	props, _ := app.listarPropriedadesSelect()
	data := map[string]interface{}{
		"Analise":      a,
		"Propriedades": props,
		"Title":        title,
	}
	app.renderTemplate(w, r, "analises/editar_sidebar.html", data)
}
 
func (app *Application) SalvarAnalise(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.serverError(w, r, err)
		return
	}
	id := r.Form.Get("id")
	propID, _ := strconv.Atoi(r.Form.Get("propriedade_id"))
	talhao := r.Form.Get("talhao")
	dataAmostra := r.Form.Get("data_amostra")
	ph, _ := strconv.ParseFloat(r.Form.Get("ph"), 64)
	mo, _ := strconv.ParseFloat(r.Form.Get("materia_organica"), 64)
	n, _ := strconv.ParseFloat(r.Form.Get("nitrogenio"), 64)
	p, _ := strconv.ParseFloat(r.Form.Get("fosforo"), 64)
	k, _ := strconv.ParseFloat(r.Form.Get("potassio"), 64)
	ctc, _ := strconv.ParseFloat(r.Form.Get("ctc"), 64)
	sb, _ := strconv.ParseFloat(r.Form.Get("saturacao_bases"), 64)
	recCalc := r.Form.Get("rec_calcario")
	recAdub := r.Form.Get("rec_adubacao")
	plano := r.Form.Get("plano_correcao")
	status := r.Form.Get("status")
	if status == "" {
		status = "pendente"
	}
 
	if id == "" || id == "0" {
		_, err := app.DB.Exec(`INSERT INTO analises (propriedade_id, talhao, data_amostra, ph, materia_organica, nitrogenio, fosforo, potassio, ctc, saturacao_bases, rec_calcario, rec_adubacao, plano_correcao, status) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			propID, talhao, dataAmostra, ph, mo, n, p, k, ctc, sb, recCalc, recAdub, plano, status)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else {
		idInt, _ := strconv.Atoi(id)
		_, err := app.DB.Exec(`UPDATE analises SET propriedade_id=?, talhao=?, data_amostra=?, ph=?, materia_organica=?, nitrogenio=?, fosforo=?, potassio=?, ctc=?, saturacao_bases=?, rec_calcario=?, rec_adubacao=?, plano_correcao=?, status=? WHERE id=?`,
			propID, talhao, dataAmostra, ph, mo, n, p, k, ctc, sb, recCalc, recAdub, plano, status, idInt)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	w.Header().Set("HX-Redirect", "/analises")
	w.WriteHeader(http.StatusOK)
}
 
func (app *Application) DetalhesAnalise(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	var a AnaliseSolo
	row := app.DB.QueryRow(`SELECT a.id, a.propriedade_id, p.nome, c.nome, a.talhao, a.data_amostra, a.ph, a.materia_organica, a.nitrogenio, a.fosforo, a.potassio, a.ctc, a.saturacao_bases, a.rec_calcario, a.rec_adubacao, a.plano_correcao, a.status, a.data_cadastro
		FROM analises a JOIN propriedades p ON p.id = a.propriedade_id JOIN clientes c ON c.id = p.cliente_id WHERE a.id = ?`, id)
	err = row.Scan(&a.ID, &a.PropriedadeID, &a.PropriedadeNome, &a.ClienteNome, &a.Talhao, &a.DataAmostra, &a.PH, &a.MateriaOrganica, &a.Nitrogenio, &a.Fosforo, &a.Potassio, &a.CTC, &a.SaturacaoBases, &a.RecCalcario, &a.RecAdubacao, &a.PlanoCorrecao, &a.Status, &a.DataCadastro)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
	data := map[string]interface{}{
		"Analise": a,
		"Title":   "Detalhes da Análise",
	}
	app.renderTemplate(w, r, "analises/detalhes_sidebar.html", data)
}
 
func (app *Application) ExcluirAnalise(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	app.DB.Exec("DELETE FROM analises WHERE id = ?", id)
	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Análise excluída com sucesso.", "type": "success"}}`)
	w.WriteHeader(http.StatusOK)
}
 
func (app *Application) listarPropriedadesSelect() ([]Propriedade, error) {
	rows, err := app.DB.Query(`SELECT p.id, p.nome, c.nome FROM propriedades p JOIN clientes c ON c.id = p.cliente_id ORDER BY p.nome`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var props []Propriedade
	for rows.Next() {
		var p Propriedade
		rows.Scan(&p.ID, &p.Nome, &p.ClienteNome)
		props = append(props, p)
	}
	return props, nil
}

