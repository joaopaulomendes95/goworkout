package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/strangecousinwst/goworkout/cmd/web"
	"github.com/strangecousinwst/goworkout/internal/store"
	"github.com/strangecousinwst/goworkout/internal/tokens"
)

// SignupWebHandler handles both GET and POST requests for /signup
func SignupWebHandler(userStore store.UserStore, tokenStore store.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusInternalServerError)
				return
			}

			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")
			bio := r.FormValue("bio")

			// Validate inputs - you can add more validation rules if needed
			if username == "" || email == "" || password == "" {
				web.SignupForm("All fields are required").Render(r.Context(), w)
				return
			}

			// Check if username exists
			existingUser, err := userStore.GetUserByUsername(username)
			if err != nil {
				log.Printf("Error checking username: %v", err)
				web.SignupForm("Error creating account").Render(r.Context(), w)
				return
			}

			if existingUser != nil {
				web.SignupForm("Username already taken").Render(r.Context(), w)
				return
			}

			// Create user
			user := &store.User{
				Username: username,
				Email:    email,
				Bio:      bio,
			}

			// Set password (this handles the hashing)
			err = user.PasswordHash.Set(password)
			if err != nil {
				log.Printf("Error hashing password: %v", err)
				web.SignupForm("Error creating account").Render(r.Context(), w)
				return
			}

			// Insert user into database
			err = userStore.CreateUser(user)
			if err != nil {
				log.Printf("Error creating user: %v", err)
				web.SignupForm("Error creating account").Render(r.Context(), w)
				return
			}

			// Create authentication token
			token, err := tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
			if err != nil {
				log.Printf("Error creating token: %v", err)
				// User created but token failed - redirect to login
				http.Redirect(w, r, "/login?msg=account_created", http.StatusSeeOther)
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

		// GET - show signup form
		err := web.SignupForm("").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering signup form: %v", err)
		}
	}
}
