package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/bbquite/go-loyalty/internal/utils"
)

func TestMW(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		token := strings.Split(authHeader, " ")

		if len(token) == 2 {

			authToken := token[1]
			authorized, err := utils.IsAuthorized(authToken)
			if err != nil {
				log.Print(err)
			}

			if authorized {
				userID, err := utils.ExtractIDFromToken(authToken)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				r.Header.Set("x-user-id", userID)
				return
			}
		}

		w.WriteHeader(http.StatusUnauthorized)
	}
}
