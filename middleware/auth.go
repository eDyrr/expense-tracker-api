package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	key   []byte
	Store *sessions.CookieStore
)

func init() {
	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the session key
	key = []byte(os.Getenv("KEY"))
	if len(key) == 0 {
		log.Fatal("KEY environment variable is not set")
	}

	// Initialize the session store
	Store = sessions.NewCookieStore(key)

	// Set session options
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
	}
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "authentification")

		fmt.Printf("Session values: %v\n", session.Values) // Debug session values

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
