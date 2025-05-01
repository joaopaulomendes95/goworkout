package handlers

import (
	"github.com/strangecousinwst/goworkout/internal/store"
)

// WebHandlers holds dependencies for web route handlers
type WebHandlers struct {
	workoutStore store.WorkoutStore
	userStore    store.UserStore
	tokenStore   store.TokenStore
}

// NewWebHandlers creates a new WebHandlers instance
func NewWebHandlers(
	workoutStore store.WorkoutStore,
	userStore store.UserStore,
	tokenStore store.TokenStore,
) *WebHandlers {
	return &WebHandlers{
		workoutStore: workoutStore,
		userStore:    userStore,
		tokenStore:   tokenStore,
	}
}
