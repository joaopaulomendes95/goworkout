package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/strangecousinwst/goworkout/cmd/web"
	"github.com/strangecousinwst/goworkout/internal/middleware"
	"github.com/strangecousinwst/goworkout/internal/store"
)

// WorkoutDetailHandler renders a single workout's details
func WorkoutDetailHandler(workoutStore store.WorkoutStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid workout ID", http.StatusBadRequest)
			return
		}

		workout, err := workoutStore.GetWorkoutByID(id)
		if err != nil {
			log.Printf("Error getting workout: %v", err)
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}

		// Check if user is authorized to view this workout
		user := middleware.GetUser(r)
		if user == nil || user == store.AnonymousUser || int64(user.ID) != int64(workout.UserID) {
			http.Error(w, "Not authorized to view this workout", http.StatusForbidden)
			return
		}

		err = web.WorkoutDetail(*workout).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering workout detail: %v", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	}
}

// WorkoutListHandler renders a list of all user workouts
func WorkoutListHandler(workoutStore store.WorkoutStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.GetUser(r)
		if user == nil || user == store.AnonymousUser {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		workouts, err := workoutStore.GetWorkoutsForUser(user.ID)
		if err != nil {
			log.Printf("Error getting workouts: %v", err)
			http.Error(w, "Error loading workouts", http.StatusInternalServerError)
			return
		}

		err = web.WorkoutList(workouts).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering workout list: %v", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	}
}

// WorkoutCreateHandler renders the form for creating a new workout
func WorkoutCreateHandler(workoutStore store.WorkoutStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.GetUser(r)
		if user == nil || user == store.AnonymousUser {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Create an empty workout for the form
		workout := &store.Workout{}

		err := web.WorkoutForm(workout).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering workout form: %v", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
		}
	}
}

// WorkoutEditHandler renders the form for editing an existing workout
func WorkoutEditHandler(workoutStore store.WorkoutStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid workout ID", http.StatusBadRequest)
			return
		}

		workout, err := workoutStore.GetWorkoutByID(id)
		if err != nil {
			log.Printf("Error getting workout: %v", err)
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}

		// Check if user is authorized to edit this workout
		user := middleware.GetUser(r)
		if user == nil || user == store.AnonymousUser || int64(user.ID) != int64(workout.UserID) {
			http.Error(w, "Not authorized to edit this workout", http.StatusForbidden)
			return
		}

		err = web.WorkoutForm(workout).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering workout form: %v", err)
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			err = web.WorkoutForm(workout).Render(r.Context(), w)
			if err != nil {
				log.Printf("Error rendering workout form: %v", err)
				http.Error(w, "Error rendering page", http.StatusInternalServerError)
			}
		}
	}
}
