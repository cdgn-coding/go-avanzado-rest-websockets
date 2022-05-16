package middlewares

import (
	server "go-rest-websockets/server"
	"net/http"
	"strings"
)

var (
	NoAuthNeeded = []string{"login", "signup", "ws"}
)

func isAuthNeeded(uri string) bool {
	for _, noNeededRoute := range NoAuthNeeded {
		if strings.Contains(uri, noNeededRoute) {
			return false
		}
	}

	return true
}

func CheckAuthMiddleware(s server.Server, auth server.Authorization) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isAuthNeeded(r.RequestURI) {
				next.ServeHTTP(w, r)
				return
			}
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			_, err := auth.ParseAndVerifyToken(s.Config().JWTSecret, tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
