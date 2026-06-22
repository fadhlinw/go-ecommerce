package handler

import (
	"net/http"
)

// In a real application, these would verify a JWT token.
// For this demonstration, we'll look at a simple header: "Authorization: Admin" or "Authorization: User"

func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Admin" {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Admin" && auth != "User" {
			http.Error(w, "Unauthorized: Please log in", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
