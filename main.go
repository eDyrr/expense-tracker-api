package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/routes"
	"github.com/joho/godotenv"

	"github.com/gorilla/handlers"
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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	_, err = database.ConnectDB()
	if err != nil {
		panic("could not connect to db")
	}

	fmt.Println("Connection is successful")

	router := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Autherization", "Accept", "Origin",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Access-Control-Allow-Methods",
			"Access-Control-Expose-Headers",
			"Access-Control-Max-Age",
			"Access-Control-Allow-Credentials"}),
		handlers.AllowCredentials(),
	)

	routes.SetUpRoutes(router)

	http.ListenAndServe(":3000", cors(router))
}
