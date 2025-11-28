package middleware

import (
	"net/http"

	"platform.zone01.gr/git/gpapadopoulos/forum/config"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if the origin matches our allowed frontend
		if origin == config.Config.AllowedOrigins {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Always set these headers for API functionality
		w.Header().Set("Access-Control-Allow-Methods", config.Config.AllowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", config.Config.AllowedHeaders)
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight OPTIONS requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}