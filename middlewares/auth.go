package middlewares

import (
	r "blog-api/repositories/user"
	"fmt"
	"net/http"
	"strings"
)

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

		err := r.VerifyJWT(token)

		fmt.Println("verification error? ", err)
	}
}
