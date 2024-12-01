package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received a registration request")

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	var existingUser models.User

	if err := database.DB.Where("email = ?", data["email"]).First(&existingUser).Error; err == nil {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
	}

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: hashedPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
	}

	response := map[string]string{
		"message": "user registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received a login request")

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
	}

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0 {
		fmt.Print("user not found")
		http.Error(w, "invalid cred", http.StatusUnauthorized)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		fmt.Println("invalid password", err)
		http.Error(w, "invalid cred", http.StatusUnauthorized)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte("K+X9tyjLWZTf5gNwMTM1+jzIpBPu1Z6K2mjFZ8FmjXs="))
	if err != nil {
		fmt.Println("error generating a token", err)
		http.Error(w, "failed to generate a token", http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	response := map[string]string{
		"message": "login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func User(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request to get user")

	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "unauthorized: no token provided", http.StatusUnauthorized)
	}

	token, err := jwt.ParseWithClaims(cookie.Value, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte("K+X9tyjLWZTf5gNwMTM1+jzIpBPu1Z6K2mjFZ8FmjXs="), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "invalid token payload: missing 'sub' claim", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(sub)
	if err != nil {
		http.Error(w, "invalid user ID format", http.StatusInternalServerError)
		return
	}

	var user models.User
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
