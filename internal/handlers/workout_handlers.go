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

// WorkoutCreateHandler handles GET/POST for workout creation
func WorkoutCreateHandler(workoutStore store.WorkoutStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := middleware.GetUser(r)
		if user == nil || user == store.AnonymousUser {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodPost {
			// Process form submission
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			// Extract form values
			title := r.FormValue("title")
			description := r.FormValue("description")
			durationStr := r.FormValue("duration_minutes")
			caloriesStr := r.FormValue("calories_burned")

			// Validate required fields
			if title == "" {
				// Return form with error
				web.WorkoutForm(nil).Render(r.Context(), w)
				return
			}

			// Convert numeric values
			duration, _ := strconv.Atoi(durationStr)
			calories, _ := strconv.Atoi(caloriesStr)

			// Create workout
			workout := &store.Workout{
				UserID:          user.ID,
				Title:           title,
				Description:     description,
				DurationMinutes: duration,
				CaloriesBurned:  calories,
			}

			// Handle exercise entries if needed
			// This would parse the form entries for exercises

			// Save to database
			createdWorkout, err := workoutStore.CreateWorkout(workout)
			if err != nil {
				log.Printf("Error creating workout: %v", err)
				http.Error(w, "Could not create workout", http.StatusInternalServerError)
				return
			}

			// Redirect to the new workout
			http.Redirect(w, r, "/workouts/"+strconv.Itoa(createdWorkout.ID), http.StatusSeeOther)
			return
		}

		// GET request - show form
		err := web.WorkoutForm(nil).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering workout form: %v", err)
		}
	}
}

// WorkoutEditHandler handles viewing/editing of existing workouts
func WorkoutEditHandler(workoutStore store.WorkoutStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user and ensure authenticated
		user := middleware.GetUser(r)
		if user == nil || user == store.AnonymousUser {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get workout ID from URL
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid workout ID", http.StatusBadRequest)
			return
		}

		// Fetch workout from database
		workout, err := workoutStore.GetWorkoutByID(id)
		if err != nil {
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}

		// Ensure user owns this workout
		if workout.UserID != user.ID {
			http.Error(w, "Not authorized", http.StatusForbidden)
			return
		}

		if r.Method == http.MethodPost {
			// Process form submission for edit
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			// Update workout fields
			workout.Title = r.FormValue("title")
			workout.Description = r.FormValue("description")
			durationStr := r.FormValue("duration_minutes")
			caloriesStr := r.FormValue("calories_burned")

			// Convert numeric values
			workout.DurationMinutes, _ = strconv.Atoi(durationStr)
			workout.CaloriesBurned, _ = strconv.Atoi(caloriesStr)

			// Update in database
			err = workoutStore.UpdateWorkout(workout)
			if err != nil {
				log.Printf("Error updating workout: %v", err)
				http.Error(w, "Could not update workout", http.StatusInternalServerError)
				return
			}

			// Redirect to workout detail
			http.Redirect(w, r, "/workouts/"+idParam, http.StatusSeeOther)
			return
		}

		// GET request - show edit form with workout data
		err = web.WorkoutForm(workout).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering workout form: %v", err)
		}
	}
}

// WorkoutListHandler shows all workouts for the current user
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
		}
	}
}

// Add this handler for viewing workout details
func WorkoutDetailHandler(workoutStore store.WorkoutStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user
		user := middleware.GetUser(r)
		if user == nil || user == store.AnonymousUser {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get workout ID from URL
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Printf("Invalid workout ID: %v", err)
			http.Error(w, "Invalid workout ID", http.StatusBadRequest)
			return
		}

		// Fetch workout from database
		workout, err := workoutStore.GetWorkoutByID(id)
		if err != nil {
			log.Printf("Error fetching workout: %v", err)
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}

		// Ensure user owns this workout
		if workout.UserID != user.ID {
			http.Error(w, "Not authorized", http.StatusForbidden)
			return
		}

		// Render the workout detail template
		err = web.WorkoutDetail(*workout).Render(r.Context(), w)
		if err != nil {
			log.Printf("Error rendering workout detail: %v", err)
		}
	}
}
