package routes

import (
	"html/template"

	"github.com/eDyrr/expense-tracker-api/controllers"
	"github.com/eDyrr/expense-tracker-api/middleware"
	"github.com/eDyrr/expense-tracker-api/templates"
	"github.com/gorilla/mux"
)

func SetUpRoutes(router *mux.Router, tmpl *template.Template) {

	// API routes
	// apiRouter := router.PathPrefix("/api").Subrouter()
	router.HandleFunc("/", templates.WithTemplate(tmpl, templates.Index))
	router.HandleFunc("/signup", templates.WithTemplate(tmpl, templates.SignUp))
	router.HandleFunc("/login", templates.WithTemplate(tmpl, templates.Login))
	// apiRouter.HandleFunc("/login", controllers.Login)
	// apiRouter.HandleFunc("/auth", controllers.WithTemplate(tmpl, controllers.SignUp))

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/signup", controllers.SignUp)
	api.HandleFunc("/login", controllers.Login)

	homeMiddleware := router.PathPrefix("/site").Subrouter()
	homeMiddleware.Use(middleware.Auth)

	homeMiddleware.HandleFunc("/logout", controllers.Logout)
	homeMiddleware.HandleFunc("/home", templates.WithTemplate(tmpl, templates.Home))
	homeMiddleware.HandleFunc("/list", controllers.Listall)
	homeMiddleware.HandleFunc("/purchase", controllers.AddPurchase)
	homeMiddleware.HandleFunc("/filter", controllers.FilterPurchases)
	homeMiddleware.HandleFunc("/clicked", controllers.Select)
	homeMiddleware.HandleFunc("/delete/{id}", controllers.Delete).Methods("DELETE")
	homeMiddleware.HandleFunc("/edit/{id}", controllers.Edit)
	homeMiddleware.HandleFunc("/reset/{id}", controllers.Reset)
	homeMiddleware.HandleFunc("/submit/{id}", controllers.Submit)
}
