package utils

import (
	"embed"
	"html/template"
)

func LoadTemplates(fs embed.FS) (*template.Template, error) {
	funcmap := template.FuncMap{
							"add": func(a,b int) int { return a + b},
							"sub": func(a,b int) int { return a - b},
							"mul": func(a,b int) int { return a * b},
							"div": func(a,b int) int { return a / b},
							"formatDate" : func(date string) string {
								if date == "" {
									return ""
								}
								return date
							},
							"formatCurrency": func(value float64) string {
								return "R$ "
							},
							"formatArea": func(area float64) string {
								return " ha"
							},
	 }
	 tmpl := template.New("").Funcs(funcmap)

	 tmpl, err := tmpl.ParseFS(fs, "front-end/templates/**/*.html")
	 if err != nil {

		 return nil, err
	 }
	 return tmpl, nil
}
