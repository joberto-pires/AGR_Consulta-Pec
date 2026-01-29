package handlers

import (
	"database/sql"
	"embed"
	"html/template"
	"log"
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
	 mux.Handle("/static/", http.FileServer(http.FS(app.StaticFS)))
 }


 // Rotas Principais
 mux.HandleFunc("/", app.Homepage)
// mux.HandleFunc("/clientes", GetClients)

 return app.logRequest(mux)
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		if r.URL.Path != "/health" && r.URL.Path != "/static/" {
			log.Printf("%s %s %s", r.Method, r.URL.Path, duration)
		}
	})
}

func (app *Application) Homepage(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "Index.html", nil)

}

func (app *Application) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	// Adicionar dados comuns a todos os templates
	templateData := map[string]interface{}{
		"Data":       data,
		"Env":        app.Env,
		"CurrentURL": r.URL.Path,
		"Year":       time.Now().Year(),
	}
	
	// Mesclar com dados específicos se for um map
	if dataMap, ok := data.(map[string]interface{}); ok {
		for k, v := range dataMap {
			templateData[k] = v
		}
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := app.Template.ExecuteTemplate(w, name, templateData); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌ Server Error: %s %s - %v", r.Method, r.URL.Path, err)
	
	if app.Env == "development" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
