package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/daniarmas/chat/config"
	"github.com/daniarmas/chat/pkg/jwt_utils"
	"github.com/google/uuid"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// A stand-in for our database backed user object
type UserContext struct {
	ID uuid.UUID
}

// Middleware decodes the share session cookie and packs the session into context
func AuthorizationMiddleware(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			t := strings.Split(authHeader, " ")

			if len(t) == 2 {
				authToken := t[1]
				authorized, err := jwt_utils.IsAuthorized(authToken, cfg.JwtSecret)
				if authorized {
					accessTokenClaim, err := jwt_utils.ExtractTokenClaim(authToken, cfg.JwtSecret)
					if err != nil {
						http.Error(w, "The access token is invalid. Please obtain a new access token and try again.", http.StatusUnauthorized)
						return
					}
					// put it in context
					ctx := context.WithValue(r.Context(), userCtxKey, &UserContext{ID: accessTokenClaim.UserId})

					// and call the next with our new context
					r = r.WithContext(ctx)
					next.ServeHTTP(w, r)
				} else {
					switch err.Error() {
					case "Token is expired":
						http.Error(w, "The access token has expired. Please obtain a new access token and try again.", http.StatusUnauthorized)
						return
					case "signature is invalid":
						http.Error(w, "The access token is invalid. Please obtain a new access token and try again.", http.StatusUnauthorized)
						return
					default:
						http.Error(w, "The server has an internal error.", http.StatusUnauthorized)
						return
					}
				}
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *UserContext {
	raw, _ := ctx.Value(userCtxKey).(*UserContext)
	return raw
}
