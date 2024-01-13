package middlewares

import (
	"go-jwt-api/config"
	"go-jwt-api/helpers"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func JWTmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]string{"message": "Unauthorized"}
				helpers.ResponseJSON(w, http.StatusUnauthorized, response)
				return
			} else {
				// Log the error for debugging
				log.Println("Error retrieving cookie: ", err)
				response := map[string]string{"message": "Internal Server Error"}
				helpers.ResponseJSON(w, http.StatusInternalServerError, response)
				return
			}
		}

		// Retrieve Token Value
		tokenString := c.Value

		// Check Token
		claims := &config.JWTClaim{}
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})
		if err != nil {
			v, ok := err.(*jwt.ValidationError)
			if ok {
				switch v.Errors {
				case jwt.ValidationErrorClaimsInvalid:
					// Invalid Token
					response := map[string]string{"message": "Unauthorized"}
					helpers.ResponseJSON(w, http.StatusUnauthorized, response)
					return
				case jwt.ValidationErrorExpired:
					// Token expired
					response := map[string]string{"message": "Token Expired"}
					helpers.ResponseJSON(w, http.StatusUnauthorized, response)
					return
				default:
					response := map[string]string{"message": "Unauthorized"}
					helpers.ResponseJSON(w, http.StatusUnauthorized, response)
					return
				}
			} else {
				// Log the error for debugging
				log.Println("Error parsing token: ", err)
				response := map[string]string{"message": "Internal Server Error"}
				helpers.ResponseJSON(w, http.StatusInternalServerError, response)
				return
			}
		}

		if !tkn.Valid {
			// Token not valid
			response := map[string]string{"message": "Unauthorized"}
			helpers.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		}

		// If Token is Valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
