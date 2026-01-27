package handlers

import (
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

}

