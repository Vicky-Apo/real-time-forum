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

	// Create reverse proxy with WebSocket support
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// Configure the proxy to handle WebSocket upgrades
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = backendURL.Scheme
		req.URL.Host = backendURL.Host
		req.Host = backendURL.Host

		// Preserve original headers for WebSocket upgrade
		if req.Header.Get("Upgrade") == "websocket" {
			log.Printf("WebSocket upgrade request detected")
		}
	}

	// Create file server for static files
	fs := http.FileServer(http.Dir("."))

	// Main handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Proxy API, WebSocket, and uploads requests to backend
		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/ws" || strings.HasPrefix(r.URL.Path, "/uploads/") {
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
