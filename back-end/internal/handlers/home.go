package handlers

import (
	"database/sql"
	"embed"
	"html/template"
	"net/http"
)

var Tmpl *template.Template

var DB *sql.DB

func Init(templates embed.FS, db *sql.DB) {
 Tmpl = template.Must(template.New("").ParseFS(templates, "front-end/templates/**/*.html"))
 DB = db
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	 Tmpl.ExecuteTemplate(w, "Base.html", nil)
}
