package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/middleware"
	"github.com/eDyrr/expense-tracker-api/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "signup", nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
	var body struct {
		Name     string
		Email    string
		Password string
	}

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "couldnt decode the fucking body req", http.StatusBadRequest)
		return
	}

	fmt.Println("body %v", body)

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing the password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
	}

	result := database.DB.Create(&user)

	if result.Error != nil {
		http.Error(w, "failed to craete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	user.Password = ""
	json.NewEncoder(w).Encode(&user)
}

func Login(w http.ResponseWriter, r *http.Request) {

	session, _ := middleware.Store.Get(r, "authentification")
	var body struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "couldnt decode req body", http.StatusBadRequest)
		return
	}

	var user models.User

	result := database.DB.Where("email = ?", body.Email).Find(&user)
	if result.Error != nil {
		http.Error(w, "couldnt find user with this email", http.StatusInternalServerError)
		return
	}

	fmt.Println("user %v", user)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)

}

func Listall(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	result := database.DB.Find(&users)

	if result.Error != nil {
		http.Error(w, "cant load all users", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&users)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "authentification")

	session.Values["authenticated"] = false
	session.Save(r, w)
}

func Home(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.ExecuteTemplate(w, "home", nil)
	if err != nil {
		fmt.Println("not rendering the html")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func WithTemplate(tmpl *template.Template, handler func(http.ResponseWriter, *http.Request, *template.Template)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, tmpl)
	}
}

func AddPurchase(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string  `json:"Name"`
		Category string  `json:"Category"`
		Price    float32 `json:"Price"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	fmt.Print("user %v", body)

	// session, _ := middleware.Store.Get(r, "authentification")

	// purchase := models.Purchase{
	// 	Name:     body.Name,
	// 	Category: body.Category,
	// 	Cost:     body.Price,
	// 	UserID:   session.Values["user_id"].(uint),
	// }

	// result := database.DB.Create(&purchase)
	// if result.Error != nil {
	// 	http.Error(w, "", http.StatusInternalServerError)
	// 	return
	// }

	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(purchase)
}
