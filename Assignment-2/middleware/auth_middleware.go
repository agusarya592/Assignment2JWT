package middleware

import (
	"assignment2/config"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

func AuthMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[Auth Middleware] Request to %v\n", r.URL)
			token := r.Header.Get("Authorization")
			if token == "" {
				log.Printf("No token on authorization")
			}

			splittedToken := strings.Split(token, " ")
			if len(splittedToken) != 2 {
				log.Printf("Invalid token")
			}
			if splittedToken[0] != "Bearer" {
				log.Printf("Bearer not found")
			}

			accessToken := splittedToken[1]

			_, err := validateToken(accessToken)
			if err != nil {
				log.Printf("User unauthenticated")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func validateToken(tokenString string) (interface{}, error) {
	cfg := config.GetConfig()
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("sign in to proceed")
		} else if method != cfg.JWT_SIGNING_METHOD {
			return nil, errors.New("sign in to proceed")
		}
		return config.GetConfig().JWT_SECRET_KEY, nil
	})
	return token.Claims.(jwt.MapClaims), nil
}
