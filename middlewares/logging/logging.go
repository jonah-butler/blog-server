package loggingmiddleware

import (
	"log"
	"net/http"
	"time"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("[%s] %s %s from %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v\n", time.Since(start))
	})
}
