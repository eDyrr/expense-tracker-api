package routes

import (
	"github.com/eDyrr/expense-tracker-api/controllers"
	"github.com/gorilla/mux"
)

func SetUpRoutes(router *mux.Router) {
	router.HandleFunc("/", controllers.Hello)
}
