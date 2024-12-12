package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte("someShit"))

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
		Secure:   true,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

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

func Home(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {

	err := tmpl.ExecuteTemplate(w, "index.html", nil)
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
