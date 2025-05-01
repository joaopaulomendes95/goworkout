package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/strangecousinwst/goworkout/cmd/web"
)

// LoginWebHandler handles both GET and POST requests for /login
func LoginWebHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusInternalServerError)
			return
		}
		username := r.FormValue("username")
		password := r.FormValue("password")

		// TODO: Replace with your actual user store authentication
		// Example with real user store:
		// user, err := userStore.Authenticate(username, password)
		// if err == nil && user != nil {

		// For now, keeping simple auth for testing:
		if username == "admin" && password == "admin" {
			// Set a session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "user-session-token", // Replace with real token in production
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			})

			// Redirect to dashboard
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Failed login
		errorMsg := "Invalid username or password"
		err = web.LoginForm(errorMsg).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering login form: %v", err)
		}
		return
	}

	// GET request - show login form
	err := web.LoginForm("").Render(r.Context(), w)
	if err != nil {
		log.Printf("Error rendering login form: %v", err)
	}
}
