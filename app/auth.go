package app

import (
	"../models"
	u "../utils"
	"context"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

var notAuthenticated = []string{"/api/user/new", "/api/user/login"}

//WTF is this, some Szymon might ask. This function creates a middleware, which is registered to chain in main.go ->
//it is similar to what FilterChains do in spring framework, but it dosc not feature the fucked up aspect precendece of Spring.
//Generally it handles the authentication and JWT token generation
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

		if tokenHeader == "" {
			respondTokenMissing(&response, w)
			return
		}

		tokenHeaderParts := strings.Split(tokenHeader, " ")

		if len(tokenHeaderParts) != 2 {
			respondInvalidToken(&response, w)
			return
		}

		tokenPart := tokenHeaderParts[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil || !token.Valid {
			respondInvalidToken(&response, w)
			return
		}

		fmt.Sprintf("User %", tk.UserId)
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func respondInvalidToken(response *map[string]interface{}, w http.ResponseWriter) {
	*response = u.Message(false, "Invalid/Malformed auth token. Is is in format Authorization Bearer {token}?")
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	u.Respond(w, *response)
}

func respondTokenMissing(response *map[string]interface{}, w http.ResponseWriter) {
	*response = u.Message(false, "Missing auth token")
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	u.Respond(w, *response)
}
