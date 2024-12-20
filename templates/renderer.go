package templates

import (
	"html/template"
	"log"
	"net/http"

	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/middleware"
	"github.com/eDyrr/expense-tracker-api/models"
	"github.com/eDyrr/expense-tracker-api/services"
)

func WithTemplate(tmpl *template.Template, handler func(http.ResponseWriter, *http.Request, *template.Template)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, tmpl)
	}
}

func Home(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	session, _ := middleware.Store.Get(r, "authentification")

	ID, ok := session.Values["user_id"].(uint)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusUnauthorized)
		return
	}

	var user models.User
	err := database.DB.Where("id = ?", ID).First(&user).Error
	if err != nil {
		http.Error(w, "Error loading user data", http.StatusInternalServerError)
		return
	}

	user.Purchases = services.GetPurchases(ID)

	// fmt.Print(user.Purchases)

	// Handle template execution errors gracefully
	if err = tmpl.ExecuteTemplate(w, "home", user); err != nil {
		// Log the error, but don't double-call WriteHeader
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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

func Purchases(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	session, _ := middleware.Store.Get(r, "authentification")
	id := session.Values["user_id"]
	purchases := services.GetPurchases(id)
	err := tmpl.ExecuteTemplate(w, "purchases", purchases)
	if err != nil {
		http.Error(w, "error rendering the purchases template", http.StatusInternalServerError)
		return
	}
}
