package middlewares

import (
	ck "blog-api/contextkeys"
	r "blog-api/repositories/user"
	u "blog-api/utilities"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func BearerAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			error := fmt.Errorf("authorization header is missing")
			log.Println(error)
			u.WriteJSONErr(w, http.StatusBadRequest, error)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			error := fmt.Errorf("invalid token format")
			log.Println(error)
			u.WriteJSONErr(w, http.StatusBadRequest, error)
			return
		}

		token := parts[1]

		userID, err := r.VerifyJWT(token)
		if err != nil {
			error := fmt.Errorf("verification failed: %v", err)
			log.Println(error)
			u.WriteJSONErr(w, http.StatusUnauthorized, error)
			return
		}

		ctx := context.WithValue(req.Context(), ck.UserIDKey, userID)

		next(w, req.WithContext(ctx))
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("[%s] %s %s from %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v\n", time.Since(start))
	})
}
