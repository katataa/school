package handlers

import (
	"html/template"
	"strings"
)

var templates *template.Template

func InitializeTemplates() {
	funcMap := template.FuncMap{
		"replacePrefix": replacePrefix,
	}
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
}

func replacePrefix(path, oldPrefix, newPrefix string) string {
	if strings.HasPrefix(path, oldPrefix) {
		return newPrefix + strings.TrimPrefix(path, oldPrefix)
	}
	return path
}
