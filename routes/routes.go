package routes

import (
	"html/template"
	"net/http"

	"github.com/eDyrr/expense-tracker-api/controllers"
	"github.com/eDyrr/expense-tracker-api/middleware"
	"github.com/gorilla/mux"
)

func SetUpRoutes(router *mux.Router, tmpl *template.Template) {

	// Public routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "index", nil); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	})
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "login", nil); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	})

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/signup", controllers.SignUp)
	apiRouter.HandleFunc("/login", controllers.Login)

	homeMiddleware := router.PathPrefix("/site").Subrouter()
	homeMiddleware.Use(middleware.Auth)

	homeMiddleware.HandleFunc("/logout", controllers.Logout)
	homeMiddleware.HandleFunc("/home", controllers.WithTemplate(tmpl, controllers.Home))
	homeMiddleware.HandleFunc("/list", controllers.Listall)
	homeMiddleware.HandleFunc("/purchase", controllers.AddPurchase)
}
