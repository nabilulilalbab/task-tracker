package utils

import (
	"html/template"
	"log"
	"strings"
)

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
	log.Println("Parsing semua templates...")

	// Buat template baru, daftarkan fungsi, lalu parse file
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("templates/**/*.html")
	if err != nil {
		// Jika ada error (misal: pola salah), program akan berhenti dengan pesan jelas
		panic("Gagal mem-parsing templates: " + err.Error())
	}

	log.Println("Parsing templates selesai.")
	return tmpl
}
