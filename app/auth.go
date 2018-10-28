package app

import (
	"../dao"
	u "../utils"
	"context"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

var notAuthenticated = []string{"/user", "/user/login"}

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestPath := r.URL.Path

		for _, value := range notAuthenticated {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})

		tokenHeader := r.Header.Get("Authorization")

		tokenString := ""

		bearerPathVar := strings.Split(r.URL.RawQuery, "=")

		if tokenHeader != "" {
			tokenHeaderParts := strings.Split(tokenHeader, " ")

			if len(tokenHeaderParts) != 2 {
				respondInvalidToken(&response, w)
				return
			}

			tokenString = tokenHeaderParts[1]
		} else if len(bearerPathVar) != 0 {
			tokenString = bearerPathVar[1]
		} else {
			respondTokenMissing(&response, w)
			return
		}

		tk := &dao.Token{}

		token, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil || !token.Valid {
			respondInvalidToken(&response, w)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func respondInvalidToken(response *map[string]interface{}, w http.ResponseWriter) {
	*response = u.Message(401, "Invalid/Malformed auth token. Is is in format Authorization Bearer {token}?")
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	u.Respond(w, *response)
}

func respondTokenMissing(response *map[string]interface{}, w http.ResponseWriter) {
	*response = u.Message(401, "Missing auth token")
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	u.Respond(w, *response)
}
