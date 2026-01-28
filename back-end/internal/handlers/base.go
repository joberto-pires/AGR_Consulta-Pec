package handlers

import (
	FS "AGR_Consulta-Pec"
	"database/sql"
	"embed"
	"html/template"
	"net/http"
	"time"
)

type Application struct {
	DB				      *sql.DB
	Template      	*template.Template
	Env				 			string
	StaticFS   			embed.FS
	StartTime  			time.Time
}

func (app *Application) Routes() http.Handler {
 mux := http.NewServeMux()

 if app.Env == "development" {
	 fs := http.FileServer(http.Dir("front-end/static"))
	 mux.Handle("/static/", http.StripPrefix("/static/", fs))
 } else {
	 mux.Handle("front-end/static/", http.FileServer(http.FS(app.StaticFS)))
 }

  Init(FS.TemplatesFS, app.DB)

 // Rotas Principais
 mux.HandleFunc("/", Homepage)
 mux.HandleFunc("/clientes", GetClients)

 return mux
}

