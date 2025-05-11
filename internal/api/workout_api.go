package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/strangecousinwst/goworkout/internal/middleware"
	"github.com/strangecousinwst/goworkout/internal/store"
	"github.com/strangecousinwst/goworkout/internal/utils"
)

// WorkoutAPI holds dependencies for workout-related API handlers.
type WorkoutAPI struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

// NewWorkoutAPI creates a new WorkoutAPI instance with the provided workout store and logger.
func NewWorkoutAPI(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutAPI {
	return &WorkoutAPI{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

// HandleGetUserWorkouts handles requests to get all workouts for the authenticated user.
// It checks if the user is authenticated and retrieves their workouts from the store.
func (wh *WorkoutAPI) HandleGetUserWorkouts(w http.ResponseWriter, r *http.Request) {
	// Get the current user from the request context
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser.IsAnonymous() {
		wh.logger.Println("WARN: [WorkoutAPI.HandleGetUserWorkouts] Attempt to get workouts without authenticated user.")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Authentication required"})
		return
	}

	// Get the workouts for the current user
	wh.logger.Printf("INFO: [WorkoutAPI.HandleGetUserWorkouts] Fetching workouts for user ID: %d", currentUser.ID)
	workouts, err := wh.workoutStore.GetWorkoutsForUser(currentUser.ID)
	if err != nil {
		wh.logger.Printf("ERROR: [WorkoutAPI.HandleGetUserWorkouts] Failed to get workouts for user ID %d: %v", currentUser.ID, err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve workouts"})
		return
	}

	// If no workouts are found, initialize an empty slice
	if workouts == nil {
		workouts = []store.Workout{}
	}

	wh.logger.Printf("INFO: [WorkoutAPI.HandleGetUserWorkouts] Successfully fetched %d workouts for user ID: %d", len(workouts), currentUser.ID)
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workouts": workouts})
}

// HandleGetWorkoutByID handles requests to get a specific workout by its ID.
// It checks if the user is authenticated and retrieves the workout from the store.
func (wh *WorkoutAPI) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	// Get the id from the request context
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout ID"})
		return
	}

	// Get the workout given the ID
	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "invalid server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

// HandleCreateWorkout handles requests to create a new workout.
// It checks if the user is authenticated and creates the workout in the store.
func (wh *WorkoutAPI) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	// Decode the request body into the workout struct
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decodingCreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	// Get the current user from the request context
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
	}

	workout.UserID = currentUser.ID

	// Create the workout in the database
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: creatingWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

// HandleUpdateWorkoutByID handles requests to update a specific workout by its ID.
func (wh *WorkoutAPI) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the request body
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout update ID"})
		return
	}

	// Get the workout by ID
	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	// At this point we can assume we are able to find the workout
	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	// Decode the request body into the updateWorkoutRequest struct
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.logger.Printf("ERROR: decondingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	// Validate the request
	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	// Get the current user from the request context
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	// Check owner of the workout
	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	// Check if the current user is the owner of the workout
	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you do not have permission to update this workout"})
		return
	}

	// Update the workout in the database
	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: updatingWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

// HandleDeleteWorkoutByID handles requests to delete a specific workout by its ID.
func (wh *WorkoutAPI) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the request context
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout delete ID"})
		return
	}

	// Get the user from the request context
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
		return
	}

	// Get the workout owner
	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	// Check if the current user is the owner of the workout
	if workoutOwner != currentUser.ID {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you do not have permission to delete this workout"})
		return
	}

	// Delete the workout from the database
	err = wh.workoutStore.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}

	if err != nil {
		wh.logger.Printf("ERROR: deletingWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
