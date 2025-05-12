package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/strangecousinwst/goworkout/internal/middleware"
	"github.com/strangecousinwst/goworkout/internal/store"
	"github.com/strangecousinwst/goworkout/internal/utils"
)

// registerUserRequest defines the expected JSON structure for a user registration request
type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

// UserAPI holds dependencies for user-related API handlers.
type UserAPI struct {
	userStore store.UserStore
	logger    *log.Logger
}

// updateUserRequest defines the expected JSON structure for a user update request
type updateUserRequest struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

// NewUserAPI creates a new UserAPI instance with the provided user store and logger.
func NewUserAPI(userStore store.UserStore, logger *log.Logger) *UserAPI {
	return &UserAPI{
		userStore: userStore,
		logger:    logger,
	}
}

// ValidateRegisterRequest validates the fields of a user registration request.
func (h *UserAPI) ValidateRegisterRequest(req *registerUserRequest) error {
	// Validate username
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	} else if len(req.Username) > 20 {
		return errors.New("username must be at most 20 characters long")
	}

	// Validate email
	if req.Email == "" {
		return errors.New("email is required")
	}
	if len(req.Email) < 5 {
		return errors.New("email must be at least 5 characters long")
	} else if len(req.Email) > 50 {
		return errors.New("username must be at most 50 characters long")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	// Validate password
	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

// HandleRegisterUser handles requests for new user registration.
// It validates the input, creates a new user with a hashed password,
// and stores the user in the database.
func (h *UserAPI) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = h.ValidateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	// Set and hash the password
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	// Create the user in the database
	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: registering user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}

// HandleGetCurrentUser retrieves and returns the details of the currently authenticated user.
// The user is identified via the authentication token processed by middleware.
func (h *UserAPI) HandleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve user from context (set by authentication middleware).
	user := middleware.GetUser(r)
	if user == nil || user.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

// HandleUpdateUser handles requests to update a user's profile information.
// It allows modification of username and bio for the authenticated user.
func (h *UserAPI) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req updateUserRequest

	// Get the currently authenticated user.
	currentUser := middleware.GetUser(r)
	if currentUser.IsAnonymous() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "User is anonymous"})
		return
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding update user request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	// Update user fields
	currentUser.Username = req.Username
	currentUser.Bio = req.Bio

	err = h.userStore.UpdateUser(currentUser)
	if err != nil {
		h.logger.Printf("ERROR: updating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "username is taken"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": currentUser})
}
