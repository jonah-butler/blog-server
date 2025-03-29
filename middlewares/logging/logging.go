package loggingmiddleware

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("[%s]: %s %s from %s by %s", r.Method, r.URL.Path, r.Proto, getClientIP(r), r.UserAgent())

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v\n", time.Since(start))
	})
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}
