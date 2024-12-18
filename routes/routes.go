package routes

import (
	"html/template"

	"github.com/eDyrr/expense-tracker-api/controllers"
	"github.com/eDyrr/expense-tracker-api/middleware"
	"github.com/gorilla/mux"
)

func SetUpRoutes(router *mux.Router, tmpl *template.Template) {

	// API routes
	// apiRouter := router.PathPrefix("/api").Subrouter()
	router.HandleFunc("/", controllers.WithTemplate(tmpl, controllers.Index))
	router.HandleFunc("/signup", controllers.WithTemplate(tmpl, controllers.SignUp))
	// apiRouter.HandleFunc("/login", controllers.Login)
	// apiRouter.HandleFunc("/auth", controllers.WithTemplate(tmpl, controllers.SignUp))
	homeMiddleware := router.PathPrefix("/site").Subrouter()
	homeMiddleware.Use(middleware.Auth)

	homeMiddleware.HandleFunc("/logout", controllers.Logout)
	homeMiddleware.HandleFunc("/home", controllers.WithTemplate(tmpl, controllers.Home))
	homeMiddleware.HandleFunc("/list", controllers.Listall)
	homeMiddleware.HandleFunc("/purchase", controllers.AddPurchase)
}
