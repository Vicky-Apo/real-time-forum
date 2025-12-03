package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	PORT        = ":3000"
	BACKEND_URL = "http://localhost:8080"
)

func main() {
	// Parse backend URL for proxy
	backendURL, err := url.Parse(BACKEND_URL)
	if err != nil {
		log.Fatal("Invalid backend URL:", err)
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// Create file server for static files
	fs := http.FileServer(http.Dir("."))

	// Main handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Proxy API and WebSocket requests to backend
		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/ws" {
			log.Printf("Proxying: %s %s", r.Method, r.URL.Path)
			proxy.ServeHTTP(w, r)
			return
		}

		// Check if file exists
		path := filepath.Join(".", r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// File doesn't exist - serve index.html (SPA routing)
			log.Printf("SPA route: %s -> index.html", r.URL.Path)
			http.ServeFile(w, r, "./index.html")
			return
		}

		// Serve static file
		log.Printf("Static file: %s", r.URL.Path)
		fs.ServeHTTP(w, r)
	})

	log.Printf("🚀 Frontend server running at http://localhost%s", PORT)
	log.Printf("🔄 Proxying /api/* and /ws to %s", BACKEND_URL)
	log.Printf("📁 Serving static files from: ./client")
	log.Fatal(http.ListenAndServe(PORT, nil))
}
