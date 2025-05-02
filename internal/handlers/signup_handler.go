package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/strangecousinwst/goworkout/cmd/web"
)

// SignupWebHandler handles both GET and POST requests for /signup
func SignupWebHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusInternalServerError)
				return
			}

			// Just collect form data - let the API handle validation
			requestBody := map[string]string{
				"username": r.FormValue("username"),
				"email":    r.FormValue("email"),
				"password": r.FormValue("password"),
				"bio":      r.FormValue("bio"),
			}

			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				log.Printf("Error marshaling user data: %v", err)
				web.SignupForm("Something went wrong").Render(r.Context(), w)
				return
			}

			// Forward to your API - let it handle validation
			resp, err := http.Post(
				"http:/localhost:8080/users", // Use relative path for better maintainability
				"application/json",
				bytes.NewBuffer(jsonData),
			)

			if err != nil {
				log.Printf("API request error: %v", err)
				web.SignupForm("Error connecting to server").Render(r.Context(), w)
				return
			}
			defer resp.Body.Close()

			// Handle validation errors from API
			if resp.StatusCode == http.StatusBadRequest {
				var errorResp struct {
					Error string `json:"error"`
				}

				if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
					web.SignupForm(errorResp.Error).Render(r.Context(), w)
				} else {
					web.SignupForm("Invalid form data").Render(r.Context(), w)
				}
				return
			}

			// Other error handling (same as before)
			if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
				web.SignupForm("Failed to create user").Render(r.Context(), w)
				return
			}

			// Success - authenticate and redirect (same as before)
			// ...login logic from previous implementation...

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
