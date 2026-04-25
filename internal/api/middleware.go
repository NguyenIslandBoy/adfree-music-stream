package api

import (
	"net/http"
	"os"
	"strings"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/static") {
			next.ServeHTTP(w, r)
			return
		}

		token := os.Getenv("API_TOKEN")
		if token == "" {
			writeError(w, "server misconfigured: API_TOKEN not set", http.StatusInternalServerError)
			return
		}

		// Accept token from Authorization header OR ?token= query param
		authHeader := r.Header.Get("Authorization")
		queryToken := r.URL.Query().Get("token")

		validHeader := authHeader != "" && strings.TrimSpace(authHeader) == "Bearer "+token
		validQuery := queryToken == token

		if !validHeader && !validQuery {
			writeError(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
