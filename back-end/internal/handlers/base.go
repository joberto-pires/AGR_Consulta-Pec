package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Application struct {
	DB				      *sql.DB
	TemplatesFS     	string
	Env				 			string
	StaticFS   			string
	StartTime  			time.Time

	//cache 
	templates       *template.Template
	templatesLock   sync.RWMutex
}

func (app *Application) InitTemplates() error {
	return app.ReloadTemplates()
}

func (app *Application) ReloadTemplates() error {
	app.templatesLock.Lock()
	defer app.templatesLock.Unlock()
	
	log.Printf("📄 Carregando templates de: %s", app.TemplatesFS)
	
	// Verificar se o diretório existe
	if _, err := os.Stat(app.TemplatesFS); os.IsNotExist(err) {
		log.Printf("⚠️  Diretório de templates não encontrado: %s", app.TemplatesFS)
	}
	
	// Criar template com funções
	tmpl := template.New("").Funcs(template.FuncMap{
	"add": func(a,b int) int { return a + b},
	"sub": func(a,b int) int { return a - b},
	"mul": func(a,b int) int { return a * b},
	"div": func(a,b int) int { return a / b},
	"split": strings.Split,
	"iterate": func(start, end int) []int {
		var list []int
		for i := start; i <= end; i++ {
			list = append(list, i)
		}
		return list
	},
"now": func() time.Time {
    return time.Now()
},
"formatDate": func(format string, date time.Time) string {
    return date.Format(format)
},
"formatCurrency": func(value float64) string {
		return "R$ "
	},
	"formatArea": func(area float64) string {
		return " ha"
	},	

	"firstLetter": func(s string) string {
		if len(s) > 0 {
			return string(s[0])
		}
		return ""
	},
	
	"seq": func(start, end int) []int {
		var seq []int
		for i := start; i <= end; i++ {
			seq = append(seq, i)
		}
		return seq
	},
	
	"formatPhone": func(phone string) string {
		// Formatação de telefone brasileiro
		if len(phone) == 11 {
			return fmt.Sprintf("(%s) %s-%s", phone[:2], phone[2:7], phone[7:])
		} else if len(phone) == 10 {
			return fmt.Sprintf("(%s) %s-%s", phone[:2], phone[2:6], phone[6:])
		}
		return phone
	},
	
	"formatCPFCNPJ": func(doc string) string {
		// Formatação de CPF/CNPJ
		if len(doc) == 11 {
			return fmt.Sprintf("%s.%s.%s-%s", doc[:3], doc[3:6], doc[6:9], doc[9:])
		} else if len(doc) == 14 {
			return fmt.Sprintf("%s.%s.%s/%s-%s", doc[:2], doc[2:5], doc[5:8], doc[8:12], doc[12:])
		}
		return doc
	},
	
	"truncate": func(s string, length int) string {
		if len(s) <= length {
			return s
		}
		return s[:length] + "..."
	},
	
	"json": func(v interface{}) string {
		b, _ := json.Marshal(v)
		return string(b)
	},
})	
	// Percorrer diretório de templates
	err := filepath.WalkDir(app.TemplatesFS, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		// Ignorar diretórios
		if d.IsDir() {
			return nil
		}
		
		// Apenas arquivos .html
		if filepath.Ext(path) != ".html" {
			return nil
		}
		
		// Ler arquivo
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		// Nome relativo do template
		relPath, _ := filepath.Rel(app.TemplatesFS, path)
		templateName := filepath.ToSlash(relPath)
		
		// Parse template
		_, err = tmpl.New(templateName).Parse(string(content))
		if err != nil {
			log.Printf("⚠️  Erro ao parsear template %s: %v", templateName, err)
			return err
		}
		
		log.Printf("   ✅ %s", templateName)
		return nil
	})
	
	if err != nil {
		return err
	}
	
	app.templates = tmpl
	return nil
}

func (app *Application) Routes() http.Handler {
  mux := http.NewServeMux()
	
	// Servir arquivos estáticos
	fs := http.FileServer(http.Dir(app.StaticFS))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// Rotas da aplicação
	mux.HandleFunc("/", app.Homepage)
  mux.HandleFunc("/clientes", app.ListaClientes)
  mux.HandleFunc("/clientes/novo", app.FormCliente)
  mux.HandleFunc("/clientes/editar", app.FormCliente)
  mux.HandleFunc("/clientes/salvar", app.SalvarCliente)
  mux.HandleFunc("/clientes/detalhes", app.DetalhesCliente)
  mux.HandleFunc("/clientes/excluir", app.ExcluirCliente)	
	// Rota para recarregar templates em desenvolvimento
	if app.Env == "development" {
		mux.HandleFunc("/reload-templates", app.ReloadTemplatesHandler)
	}
	
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

func (app *Application) ReloadTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	if app.Env != "development" {
		http.Error(w, "Not available in production", http.StatusForbidden)
		return
	}
	
	if err := app.ReloadTemplates(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Write([]byte("Templates recarregados com sucesso!"))
}

func (app *Application) Homepage(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "index.html", nil)
}

func (app *Application) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
	app.templatesLock.RLock()
	tmpl := app.templates
	app.templatesLock.RUnlock()
	
	if tmpl == nil {
		http.Error(w, "Templates não inicializados", http.StatusInternalServerError)
		return
	}
	
	// Dados comuns para todos os templates
	templateData := map[string]interface{}{
		"Data":       data,
		"Env":        app.Env,
		"CurrentURL": r.URL.Path,
		"Year":       time.Now().Year(),
		"Version":    "1.0.0",
	}
	
	// Mesclar com dados específicos
	if dataMap, ok := data.(map[string]interface{}); ok {
		for k, v := range dataMap {
			templateData[k] = v
		}
	}
	
	// IMPORTANTE: Definir charset UTF-8 explicitamente
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	
	// Executar template
	err := tmpl.ExecuteTemplate(w, name, templateData)
	if err != nil {
		log.Printf("❌ Erro ao executar template %s: %v", name, err)
		
		// Fallback simples
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if app.Env == "development" {
			w.Write([]byte("Template Error: " + err.Error()))
		} else {
			w.Write([]byte("Erro ao carregar página"))
		}
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



