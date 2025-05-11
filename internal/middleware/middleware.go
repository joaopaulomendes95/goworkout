package middleware

import (
	"net/http"
	"strings"

	"github.com/strangecousinwst/goworkout/internal/store"
	"github.com/strangecousinwst/goworkout/internal/tokens"
	"github.com/strangecousinwst/goworkout/internal/utils"
	"golang.org/x/net/context"
)

// UserMiddleware holds the store for user-related middleware
type UserMiddleware struct {
	UserStore store.UserStore
}

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// UserContextKey is the key used to store and retrieve the user object in the request context.
const UserContextKey = contextKey("user")

// SetUser adds the provided user object to the request's context.
// It returns a new request with the updated context.
func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

// GetUser retrieves the user object from the request's context.
// It panics if the user is not found in the context,
func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		// this should never happen
		panic("missing user in request")
	}
	return user
}

// Authenticate is a middleware that inspects the "Authorization" header for a Bearer token.
// If a valid token is found, it retrieves the corresponding user and adds them to the request context.
// If no token or an invalid token is provided, it sets an AnonymousUser in the context.
func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Within this anonymous function
		// We can interject any incoming requests to the server

		// Add "Vary: Authorization" header to inform caches that the response may vary
		// based on the Authorization header.
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		// If the Authorization header is not present, set the AnonymousUser
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// Expecting "Bearer <TOKEN>" format
		headerParts := strings.Split(authHeader, " ")
		// If the header is not in the expected format, return an error
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		// Extract the token from the header
		token := headerParts[1]
		// Validate the token and retrieve the user associated with it
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}

		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or invalid"})
			return
		}

		// if all valid, set the user in the request context
		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// RequireUser is a middleware that ensures a non-anonymous user is authenticated.
// It must be used after the Authenticate middleware. If the user in the context
// is anonymous, it returns an unauthorized error.
func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user from the request context
		user := GetUser(r)

		// If the user is anonymous, return an unauthorized error
		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in to access this route"})
			return
		}

		// user is authenticated, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
