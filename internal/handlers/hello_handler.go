package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/strangecousinwst/goworkout/cmd/web"
)

func HelloWebHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("name")
	if name == "" {
		name = "World"
	}

	// err = component.Render(r.Context(), w)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	log.Fatalf("Error rendering in HelloWebHandler: %e", err)
	// }
	web.HelloPost(name).Render(r.Context(), w)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
