package handlers

import (
	"net/http"

	"github.com/strangecousinwst/goworkout/cmd/web"
)

// LoginWebHandler handles both GET and POST requests for /login
func LoginWebHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Dummy authentication logic
	if r.Method == http.MethodPost {
		if username == "admin" && password == "admin" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		_ = web.LoginForm().Render(r.Context(), w)
		return
	}

	_ = web.LoginForm().Render(r.Context(), w)
}
