package templates

import (
	"fmt"
	"html/template"
	"net/http"
)

func WithTemplate(tmpl *template.Template, handler func(http.ResponseWriter, *http.Request, *template.Template)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, tmpl)
	}
}

func Home(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "home", nil)
	if err != nil {
		fmt.Println("not rendering the html")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func Index(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "failed to render index.html", http.StatusInternalServerError)
		return
	}
}

func SignUp(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "signup", nil)
	if err != nil {
		http.Error(w, "failed to render the Signup page", http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "login", nil)
	if err != nil {
		http.Error(w, "failed to render the login page", http.StatusInternalServerError)
		return
	}
}
