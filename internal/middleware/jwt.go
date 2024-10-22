package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/bbquite/go-loyalty/internal/utils"
	"github.com/golang-jwt/jwt/v4"
)

func TokenAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		token := strings.Split(authHeader, "Bearer ")

		if len(token) == 2 {

			authToken := token[1]
			authorized, err := utils.IsAuthorized(authToken)

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				log.Printf("token auth middleware error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if authorized {

				userID, err := utils.ExtractIDFromToken(authToken)

				if err != nil {
					log.Printf("token auth middleware error: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				ctx := context.WithValue(r.Context(), utils.AccountIDContextKey, userID)
				h.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
	}
}
