package controllers

import (
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "failed to render index.html", http.StatusInternalServerError)
		return
	}
}
