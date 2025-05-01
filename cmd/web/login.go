package web

import "net/http"

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

		_ = LoginForm().Render(r.Context(), w)
		return
	}

	_ = LoginForm().Render(r.Context(), w)
}
