package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/eDyrr/expense-tracker-api/database"
	"github.com/eDyrr/expense-tracker-api/middleware"
	"github.com/eDyrr/expense-tracker-api/models"
	"github.com/gorilla/mux"
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

	fmt.Print("body %v", body)

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

	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("HX-Redirect", "/site/home")
	w.WriteHeader(http.StatusCreated)
	user.Password = ""
	// json.NewEncoder(w).Encode(&user)
}

func Login(w http.ResponseWriter, r *http.Request) {

	fmt.Print("login endpoint in")

	session, _ := middleware.Store.Get(r, "authentification")
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "couldnt decode req body", http.StatusBadRequest)
		return
	}

	fmt.Println(body)

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
	// json.NewEncoder(w).Encode(&user)
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
	fmt.Println("in the logout")
	session, _ := middleware.Store.Get(r, "authentification")

	session.Values["authenticated"] = false
	session.Save(r, w)
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
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
	// <p><strong>ID:</strong> {{ .Id }}</p>

	response := `
	<div class="purchase-item" id="purchase-{{ .ID }}">
		<p><strong>ID:</strong> {{ .ID }}</p>
		<p><strong>Name:</strong> {{ .Name }}</p>
		<p><strong>Category:</strong> {{ .Category }}</p>
		<p><strong>Cost:</strong> {{ .Cost }}</p>
		<button type="button" 
				hx-delete="/site/delete/{{ .ID }}"
				hx-target="closest .purchase-item"
				hx-swap="outerHTML">
			Delete
		</button>
		<button type="button" 
				hx-get="/site/edit/{{ .ID }}" 
				hx-target="#purchase-{{ .ID }}" 
				hx-swap="innerHTML">
			Edit
		</button>
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

func FilterPurchases(w http.ResponseWriter, r *http.Request) {

	fmt.Print("in the filter purchases route")

	var filter = struct {
		Filter string `json:"filter"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&filter)

	if err != nil {
		fmt.Print("fucking in the error")
		http.Error(w, "couldnt filter", http.StatusBadRequest)
		return
	}

	var purchases []models.Purchase

	now := time.Now()

	today := now.Format("2006-01-02")
	fmt.Println("today", today)

	week := now.Add(-7 * 24 * time.Hour)
	fmt.Println("week ago", week)

	month := now.Add(-30 * 24 * time.Hour)
	fmt.Println("month ago", month)

	months3 := now.Add(-90 * 24 * time.Hour)

	var date string

	switch filter.Filter {
	case "last_week":
		date = week.String()
	case "last_month":
		date = month.String()
	case "last_3_months":
		date = months3.String()
	}

	result := database.DB.Where("created_at >= ?", date).Find(&purchases)

	if result.Error != nil {
		http.Error(w, "cant load all purchases", http.StatusInternalServerError)
		return
	}

	fmt.Println(purchases)
}

func Select(w http.ResponseWriter, r *http.Request) {
	response := `
	<button type="button" 
			hx-delete="/site/delete/{{ .ID }}"
			hx-target="closest .purchase-item"
			hx-swap="outerHTML">
    		Delete
	</button>
`

	tmpl, err := template.New("delete").Parse(response)
	if err != nil {
		http.Error(w, "failed to parse template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		http.Error(w, "ID not provided", http.StatusBadRequest)

		return
	}

	fmt.Println(idStr)

	purchaseId, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	session, _ := middleware.Store.Get(r, "authentification")

	userId := session.Values["user_id"].(uint)

	result := DeletePurchase(uint(purchaseId), userId)

	if result != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(200)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in the edit")

	// Extracting the ID
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		http.Error(w, "Missing ID in request", http.StatusBadRequest)
		return
	}

	fmt.Println("ID to edit", idStr)

	purchaseID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var purchase models.Purchase
	result := database.DB.Where("ID = ?", purchaseID).First(&purchase)
	if err := result.Error; err != nil {
		http.Error(w, "couldnt find row using ID", http.StatusInternalServerError)
		return
	}

	fmt.Println(purchase)
	// Example data (replace with actual DB lookup)
	data := struct {
		ID   uint64
		Name string
		Cost float32
	}{
		ID:   uint64(purchase.ID),
		Name: purchase.Name,
		Cost: purchase.Cost,
	}

	// Corrected HTML form
	response := `
	<form id="my-form-{{ .ID }}" class="my-form" hx-ext="json-enc">
	<label for="name">Name</label>
	<input id="field1" name="name" value="{{ .Name }}">
	
	<label for="field2">Cost</label>
	<input id="field2" name="cost" value="{{ .Cost }}">
	
	<label for="category">Category</label>
	<select id="category" name="category" required>
		<option value="groceries">Groceries</option>
		<option value="leisure">Leisure</option>
		<option value="electronics">Electronics</option>
		<option value="utilities">Utilities</option>
		<option value="clothing">Clothing</option>
		<option value="health">Health</option>
		<option value="others">Others</option>
	</select>
	
	<div class="button-container">
		<button type="submit" hx-confirm="Are you sure you want to confirm editing?" hx-post="/site/submit/{{ .ID }}" hx-swap="innerHTML" hx-target="#my-form-{{ .ID }}">
		Submit
		</button>
		<button type="button" hx-get="/site/reset/{{ .ID }}" hx-swap="innerHTML" hx-target="#my-form-{{ .ID }}">
		Cancel
		</button>
	</div>
	</form>

	<style>
	.button-container {
		display: flex;
		gap: 10px;
	}
	</style>
	`

	// Parse and execute the template safely
	tmpl, err := template.New("response").Parse(response)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	// Set headers before writing response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	tmpl.Execute(w, data)
}

func Reset(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in the reset")
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "error parsing ID", http.StatusInternalServerError)
		return
	}

	fmt.Println("ID :", id)

	var purchase models.Purchase
	result := database.DB.Where("ID = ?", id).Find(&purchase)
	if result.Error != nil {
		http.Error(w, "error finding purchase in the DB", http.StatusInternalServerError)
		return
	}

	response := `
	<div class="purchase-item" id="purchase-{{ .ID}}">
		<p><strong>ID:</strong> {{ .ID }}</p>
		<p><strong>Name:</strong> {{ .Name }}</p>
		<p><strong>Category:</strong> {{ .Category }}</p>
		<p><strong>Cost:</strong> {{ .Cost }}</p>
		<button type="button" 
				hx-delete="/site/delete/{{ .ID }}"
				hx-target="closest div"
				hx-swap="outerHTML">
			Delete
		</button>
		<button type="button" hx-get="/site/edit/{{ .ID }}" hx-swap="innerHTML" hx-target="#purchase-{{ .ID }}">
			Edit
		</button>
	</div>
	`

	tmpl, err := template.New("response").Parse(response)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, purchase)
}

func Submit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in the submit")
	var body struct {
		Name     string `json:"name"`
		Cost     string `json:"cost"`
		Category string `json:"category"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "couldnt decode", http.StatusInternalServerError)
		return
	}
	cost64, _ := strconv.ParseFloat(body.Cost, 32)

	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)

	var purchase models.Purchase
	result := database.DB.First(&purchase, id64)
	if result.Error != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// session, _ := middleware.Store.Get(r, "authentification")

	purchase.Name = body.Name
	purchase.Cost = float32(cost64)
	purchase.Category = body.Category
	// purchase.UserID = userID

	database.DB.Save(&purchase)

	response := `
	<div class="purchase-item" id="purchase-{{ .ID }}">
		<p><strong>ID:</strong> {{ .ID }}</p>
		<p><strong>Name:</strong> {{ .Name }}</p>
		<p><strong>Category:</strong> {{ .Category }}</p>
		<p><strong>Cost:</strong> {{ .Cost }}</p>
		<button type="button" 
				hx-delete="/site/delete/{{ .ID }}"
				hx-target="closest .purchase-item"
				hx-swap="outerHTML">
			Delete
		</button>
		<button type="button" 
				hx-get="/site/edit/{{ .ID }}" 
				hx-target="#purchase-{{ .ID }}" 
				hx-swap="innerHTML">
			Edit
		</button>
	</div>
	`

	tmpl, err := template.New("response").Parse(response)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	fmt.Println(body)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Contet-Type", "text/html")
	err = tmpl.Execute(w, purchase)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

}
