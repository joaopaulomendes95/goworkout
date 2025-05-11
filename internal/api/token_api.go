package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/strangecousinwst/goworkout/internal/store"
	"github.com/strangecousinwst/goworkout/internal/tokens"
	"github.com/strangecousinwst/goworkout/internal/utils"
)

// Token API holds dependencies for token-related requests
type TokenAPI struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

// Defines the expected request payload for creating a token
type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Creates a new TokenAPI instance with necessary dependencies
func NewTokenAPI(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenAPI {
	return &TokenAPI{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

// HandleCreateToken handles requests to create a new token for a user
// expects username and password in the request body
// returns a token if the credentials are valid
func (h *TokenAPI) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	// Decode the request body into the createTokenRequest struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: createTokenRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	// Fetch the user by username
	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: getUserByUsername: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	// Check if passwords matches password hash
	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: PasswordHash.Matches %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if !passwordsDoMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	// Create a new token for the user
	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: Creating Token %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
