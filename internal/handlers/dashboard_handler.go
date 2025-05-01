package handlers

import (
	"log"
	"net/http"

	"github.com/strangecousinwst/goworkout/cmd/web"
	"github.com/strangecousinwst/goworkout/internal/middleware"
	"github.com/strangecousinwst/goworkout/internal/store"
)

// DashboardHandler renders the dashboard showing user's workouts
func DashboardHandler(workoutStore store.WorkoutStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.GetUser(r)
		var workouts []store.Workout
		var err error

		if user != nil && user != store.AnonymousUser {
			// Get the user's workouts
			workouts, err = workoutStore.GetWorkoutsForUser(user.ID)
			if err != nil {
				log.Printf("Error getting workouts: %v", err)
				http.Error(w, "Error loading workouts", http.StatusInternalServerError)
				return
			}
		}

		err = web.Dashboard(workouts).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering dashboard: %v", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	}
}
