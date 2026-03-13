package handler

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/joaodias/golang-svelte-fullstack/internal/domain"
	"github.com/joaodias/golang-svelte-fullstack/internal/repository"
)

type UserHandler struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	logger    *log.Logger
}

func NewUserHandler(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: decoding request: %v", err)
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	if err := validateRegisterRequest(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Bio:      req.Bio,
	}

	if err := user.PasswordHash.Set(req.Password); err != nil {
		h.logger.Printf("ERROR: hashing password: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to process password"})
		return
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		h.logger.Printf("ERROR: creating user: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to create user"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{"user": user})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: decoding request: %v", err)
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "username and password are required"})
		return
	}

	user, err := h.userRepo.GetByUsername(r.Context(), req.Username)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	matches, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: matching password: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal error"})
		return
	}

	if !matches {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	token := generateToken(user.ID)

	if err := h.tokenRepo.Create(r.Context(), token); err != nil {
		h.logger.Printf("ERROR: creating token: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to create session"})
		return
	}

	tokenString := hex.EncodeToString(token.Hash)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user.IsAnonymous() {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "not authenticated"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"user": user})
}

func validateRegisterRequest(req *RegisterRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 3 || len(req.Username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

func generateToken(userID int) *domain.Token {
	plainText := time.Now().UnixNano()
	return &domain.Token{
		UserID: userID,
		Hash:   domain.GenerateTokenHash(string(plainText)),
		Scope:  "authentication",
		Expiry: time.Now().Add(24 * time.Hour * 30), // 30 days
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
