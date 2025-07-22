package view

import (
	"embed"
	"html/template"
	"log"
	"strings"
)

//go:embed templates/**/*.html
var templateFiles embed.FS // Path ini sekarang valid!

var funcMap = template.FuncMap{
	"split": func(s, sep string) []string {
		return strings.Split(s, sep)
	},
	"trim": func(s string) string {
		return strings.TrimSpace(s)
	},
	"mod": func(i, j int) int {
		return i % j
	},
}

func ParseTemplates() *template.Template {
	log.Println("Parsing semua templates dari embed FS...")

	// Gunakan ParseFS untuk mem-parsing dari variabel embed.FS
	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFiles, "templates/**/*.html")
	if err != nil {
		panic("Gagal mem-parsing templates dari embed FS: " + err.Error())
	}

	log.Println("Parsing templates selesai.")
	return tmpl
}
