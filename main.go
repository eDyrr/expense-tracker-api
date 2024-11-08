package main

import (
	"fmt"
	"net/http"

	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/routes"

	"github.com/gorilla/mux"
)

// var tmpl *template.Template

func main() {
	// tmpl, _ = template.ParseGlob("./templates/*.html")

	// router := mux.NewRouter()

	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl.ExecuteTemplate(w, "index.html", nil)
	// })

	// router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {

	// })

	// http.ListenAndServe(":3000", router)

	_, err := database.ConnectDB()
	if err != nil {
		panic("could not connect to db")
	}

	fmt.Println("Connection is successful")

	router := mux.NewRouter()

	routes.SetUpRoutes(router)

	http.ListenAndServe(":3000", router)
}
