package routes

import (
	"html/template"
	"net/http"

	"github.com/eDyrr/expense-tracker-api/controllers"
	"github.com/gorilla/mux"
)

func SetUpRoutes(router *mux.Router, tmpl *template.Template) {

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "login.html", nil)
	})
	router.HandleFunc("/api/signup", controllers.SignUp)
	router.HandleFunc("/api/login", controllers.Login)
	router.HandleFunc("/list", controllers.Listall)
	router.HandleFunc("/home", controllers.WithTemplate(tmpl, controllers.Home))

}
