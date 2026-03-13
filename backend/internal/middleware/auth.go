package middleware

import (
	"context"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/joaodias/golang-svelte-fullstack/internal/domain"
	"github.com/joaodias/golang-svelte-fullstack/internal/repository"
)

type AuthMiddleware struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthMiddleware(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *AuthMiddleware {
	return &AuthMiddleware{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ctx := context.WithValue(r.Context(), "user", domain.AnonymousUser)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenHex := parts[1]
		tokenBytes, err := hex.DecodeString(tokenHex)
		if err != nil {
			http.Error(w, "invalid token format", http.StatusUnauthorized)
			return
		}

		token, err := m.tokenRepo.GetByHashAndScope(r.Context(), tokenBytes, "authentication")
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		user, err := m.userRepo.GetByID(r.Context(), token.UserID)
		if err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("user").(*domain.User)
		if !ok || user.IsAnonymous() {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
