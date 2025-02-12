package middlewares

import (
	r "blog-api/repositories/user"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func BearerAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		userID, err := r.VerifyJWT(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Verification failed: %v", err), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), UserIDKey, userID)

		next(w, req.WithContext(ctx))
	}
}
