package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/middleware"
	"github.com/eDyrr/expense-tracker-api/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Print("in the signup endpoint")

	var body struct {
		Name     string
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&body)
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

	fmt.Print("login endpoint in")

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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	w.Header().Set("HX-Redirect", "/site/home")
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

func AddPurchase(w http.ResponseWriter, r *http.Request) {
	fmt.Print("in the add purchase route")
	var body struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		Cost     string `json:"cost"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	fmt.Println(body.Name)
	fmt.Println(body.Category)
	fmt.Println(body.Cost)

	if err != nil {
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	session, _ := middleware.Store.Get(r, "authentification")

	cost64, _ := strconv.ParseFloat(body.Cost, 32)

	purchase := models.Purchase{
		Name:     body.Name,
		Category: body.Category,
		Cost:     float32(cost64),
		UserID:   session.Values["user_id"].(uint),
	}

	result := database.DB.Create(&purchase)
	if result.Error != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	response := `
	<div class="purchase-item">
		<p><strong>Name:</strong> {{.Name}}</p>
		<p><strong>Category:</strong> {{.Category}}</p>
		<p><strong>Cost:</strong> {{.Cost}}</p>
	</div>
	`

	tmpl, err := template.New("purchase").Parse(response)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, purchase)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
