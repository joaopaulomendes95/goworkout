package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/strangecousinwst/goworkout/cmd/web"
	"github.com/strangecousinwst/goworkout/internal/store"
)

// SignupWebHandler handles both GET and POST requests for /signup
func SignupWebHandler(userStore store.UserStore) http.HandlerFunc {
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

			// Validate inputs
			if username == "" || email == "" || password == "" {
				web.SignupForm("All fields are required").Render(r.Context(), w)
				return
			}

			// Check if username is available
			// Uncomment when using real user store:
			/*
			   _, err = userStore.GetUserByUsername(username)
			   if err == nil {
			       web.SignupForm("Username already taken").Render(r.Context(), w)
			       return
			   }

			   // Create the user
			   user := &store.User{
			       Username: username,
			       Email:    email,
			   }

			   err = userStore.CreateUser(user, password)
			   if err != nil {
			       log.Printf("Error creating user: %v", err)
			       web.SignupForm("Failed to create account").Render(r.Context(), w)
			       return
			   }
			*/

			// For demo purposes - skip actual user creation
			log.Printf("Would create user: %s, %s", username, email)

			// Set a session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "user-session-token", // Replace with real token
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			})

			// Redirect to dashboard
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// GET - show signup form
		err := web.SignupForm("").Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering signup form: %v", err)
		}
	}
}
