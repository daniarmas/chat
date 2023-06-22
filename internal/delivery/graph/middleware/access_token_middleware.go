package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/daniarmas/chat/internal/datasource/jwtds"
	"github.com/daniarmas/chat/internal/repository"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// A stand-in for our database backed user object
type UserContext struct {
	ID string
}

// Middleware decodes the share session cookie and packs the session into context
func AuthorizationMiddleware(jwtDs jwtds.JwtDatasource, accessTokenRepo repository.AccessTokenRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			t := strings.Split(authHeader, " ")

			if len(t) == 2 {
				authToken := t[1]
				authorized, err := jwtDs.IsAuthorized(authToken)
				if authorized {
					accessTokenClaim, err := jwtDs.ExtractTokenClaim(authToken)
					if err != nil {
						http.Error(w, "The access token is invalid. Please obtain a new access token and try again.", http.StatusUnauthorized)
						return
					}
					// Check if the user is logged in the system
					_, err = accessTokenRepo.GetAccessTokenById(context.Background(), accessTokenClaim.ID)
					if err != nil {
						http.Error(w, "Unauthorized.", http.StatusUnauthorized)
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
						http.Error(w, "The access token is invalid. Please obtain a new access token and try again.", http.StatusUnauthorized)
						return
					}
				}
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func AuthorizationWebsocketMiddleware(ctx context.Context, jwtDs jwtds.JwtDatasource, initPayload transport.InitPayload) (context.Context, error) {
	// Get the token from payload
	any := initPayload["authorization"]
	token, ok := any.(string)
	if !ok || token == "" {
		return nil, errors.New("access token is missing")
	}

	t := strings.Split(token, " ")

	if len(t) == 2 {
		authToken := t[1]
		authorized, err := jwtDs.IsAuthorized(authToken)
		if authorized {
			accessTokenClaim, err := jwtDs.ExtractTokenClaim(authToken)
			if err != nil {
				return nil, errors.New("access token invalid")
			}
			// put it in context
			ctx := context.WithValue(ctx, userCtxKey, &UserContext{ID: accessTokenClaim.UserId})

			return ctx, nil
		} else {
			switch err.Error() {
			case "Token is expired":
				return nil, errors.New("access token has expired")
			case "signature is invalid":
				return nil, errors.New("access token is invalid")
			default:
				return nil, errors.New("internal server error")
			}
		}
	} else {
		return nil, errors.New("access token is missing")
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *UserContext {
	raw, _ := ctx.Value(userCtxKey).(*UserContext)
	return raw
}
