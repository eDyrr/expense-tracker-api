package routes

import (
	"github.com/eDyrr/expense-tracker-api/controllers"
	"github.com/gorilla/mux"
)

func SetUpRoutes(router *mux.Router) {
	router.HandleFunc("/", controllers.Hello)
	router.HandleFunc("/api/register", controllers.Register)
	router.HandleFunc("/api/login", controllers.Login)
	router.HandleFunc("/api/user", controllers.User)
}
