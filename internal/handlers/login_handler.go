package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/strangecousinwst/goworkout/cmd/web"
	"github.com/strangecousinwst/goworkout/internal/store"
	"github.com/strangecousinwst/goworkout/internal/tokens"
)

// LoginWebHandler handles both GET and POST requests for /login
func LoginWebHandler(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusInternalServerError)
				return
			}

			username := r.FormValue("username")
			password := r.FormValue("password")

			if username == "" || password == "" {
				web.LoginForm("Username and password required").Render(r.Context(), w)
				return
			}

			// Authenticate user
			user, err := userStore.GetUserByUsername(username)
			if err != nil || user == nil {
				web.LoginForm("Invalid username or password").Render(r.Context(), w)
				return
			}

			// Check password
			match, err := user.PasswordHash.Matches(password)
			if err != nil || !match {
				web.LoginForm("Invalid username or password").Render(r.Context(), w)
				return
			}

			// Create authentication token
			token, err := tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
			if err != nil {
				log.Printf("Error creating token: %v", err)
				web.LoginForm("Authentication error").Render(r.Context(), w)
				return
			}

			// Set auth cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    token.PlainText,
				Path:     "/",
				Expires:  token.Expiry,
				HttpOnly: true,
			})

			// Redirect to dashboard
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		// GET request - show login form
		msg := r.URL.Query().Get("msg")
		var errorMsg string

		if msg == "account_created" {
			errorMsg = "Account created successfully. Please log in."
		}

		err := web.LoginForm(errorMsg).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering login form: %v", err)
		}
	}
}
