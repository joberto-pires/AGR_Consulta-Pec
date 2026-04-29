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
	"gtInt": func(a,b int) bool { return a > b },
	"geInt": func(a,b int) bool { return a >= b },
	"leInt": func(a,b int) bool { return a <= b },
	"ltInt": func(a,b int) bool { return a < b },
	"eqInt": func(a,b int) bool { return a == b },
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
    
    // Servir arquivos estáticos com tipos MIME corretos
    mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extrair a extensão do arquivo
        ext := filepath.Ext(r.URL.Path)
        
        // Definir Content-Type baseado na extensão
        switch ext {
        case ".css":
            w.Header().Set("Content-Type", "text/css; charset=utf-8")
        case ".js":
            w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
        case ".json":
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
        case ".png":
            w.Header().Set("Content-Type", "image/png")
        case ".jpg", ".jpeg":
            w.Header().Set("Content-Type", "image/jpeg")
        case ".gif":
            w.Header().Set("Content-Type", "image/gif")
        case ".svg":
            w.Header().Set("Content-Type", "image/svg+xml")
        case ".ico":
            w.Header().Set("Content-Type", "image/x-icon")
        case ".html", ".htm":
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
        default:
            w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        }
        
        // Desabilitar cache em desenvolvimento
        if app.Env == "development" {
            w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
            w.Header().Set("Pragma", "no-cache")
            w.Header().Set("Expires", "0")
        }
        
        // Servir o arquivo
        fs := http.FileServer(http.Dir(app.StaticFS))
        fs.ServeHTTP(w, r)
    })))
    
    // Dashboard
    mux.HandleFunc("/", app.Homepage)

    // Clientes
    mux.HandleFunc("/clientes", app.ListaClientes)
    mux.HandleFunc("/clientes/novo", app.FormCliente)
    mux.HandleFunc("/clientes/editar", app.FormCliente)
    mux.HandleFunc("/clientes/salvar", app.SalvarCliente)
    mux.HandleFunc("/clientes/detalhes", app.DetalhesCliente)
    mux.HandleFunc("/clientes/excluir", app.ExcluirCliente)

    // Propriedades
    mux.HandleFunc("/propriedades", app.ListaPropriedades)
    mux.HandleFunc("/propriedades/novo", app.FormPropriedade)
    mux.HandleFunc("/propriedades/editar", app.FormPropriedade)
    mux.HandleFunc("/propriedades/salvar", app.SalvarPropriedade)
    mux.HandleFunc("/propriedades/detalhes", app.DetalhesPropriedade)
    mux.HandleFunc("/propriedades/excluir", app.ExcluirPropriedade)

    // Análises de Solo
    mux.HandleFunc("/analises", app.ListaAnalises)
    mux.HandleFunc("/analises/novo", app.FormAnalise)
    mux.HandleFunc("/analises/editar", app.FormAnalise)
    mux.HandleFunc("/analises/salvar", app.SalvarAnalise)
    mux.HandleFunc("/analises/detalhes", app.DetalhesAnalise)
    mux.HandleFunc("/analises/excluir", app.ExcluirAnalise)

    // Consultas Técnicas
    mux.HandleFunc("/consultas", app.ListaConsultas)
    mux.HandleFunc("/consultas/novo", app.FormConsulta)
    mux.HandleFunc("/consultas/editar", app.FormConsulta)
    mux.HandleFunc("/consultas/salvar", app.SalvarConsulta)
    mux.HandleFunc("/consultas/detalhes", app.DetalhesConsulta)
    mux.HandleFunc("/consultas/excluir", app.ExcluirConsulta)

    // Monitoramento
    mux.HandleFunc("/monitoramento", app.ListaMonitoramento)
    mux.HandleFunc("/monitoramento/novo", app.FormMonitoramento)
    mux.HandleFunc("/monitoramento/editar", app.FormMonitoramento)
    mux.HandleFunc("/monitoramento/salvar", app.SalvarMonitoramento)
    mux.HandleFunc("/monitoramento/detalhes", app.DetalhesMonitoramento)
    mux.HandleFunc("/monitoramento/excluir", app.ExcluirMonitoramento)

    // Recarregar templates em desenvolvimento
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
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	app.renderTemplate(w, r, "index.html", nil)
}

func (app *Application) renderTemplate(w http.ResponseWriter, r *http.Request, name string, data any) {
    app.templatesLock.RLock()
    tmpl := app.templates
    app.templatesLock.RUnlock()
    
    if tmpl == nil {
        log.Printf("❌ Templates não inicializados ao tentar renderizar: %s", name)
        http.Error(w, "Templates não inicializados", http.StatusInternalServerError)
        return
    }
    
    log.Printf("📄 Tentando renderizar template: %s", name)
    
    // Verificar se o template existe
    if tmpl.Lookup(name) == nil {
        log.Printf("❌ Template não encontrado: %s", name)
        log.Printf("📋 Templates disponíveis:")
        tmpls := tmpl.Templates()
        for _, t := range tmpls {
            log.Printf("   - %s", t.Name())
        }
        http.Error(w, fmt.Sprintf("Template %s não encontrado", name), http.StatusInternalServerError)
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
    
    log.Printf("📊 Dados para template %s: %+v", name, templateData)
    
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
    } else {
        log.Printf("✅ Template %s renderizado com sucesso", name)
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
