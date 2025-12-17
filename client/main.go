package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	PORT        = getEnv("PORT", ":3000")
	BACKEND_URL = getEnv("BACKEND_URL", "http://localhost:8080")
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
			// Ensure WebSocket headers are preserved
			req.Header.Set("Connection", "Upgrade")
			req.Header.Set("Upgrade", "websocket")
		}
	}

	// Handle WebSocket upgrades properly
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Don't modify WebSocket upgrade responses
		return nil
	}

	// Create file server for static files
	fs := http.FileServer(http.Dir("."))

	// WebSocket upgrader for client connections
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
	}

	// WebSocket proxy handler
	wsProxy := func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a WebSocket upgrade request
		if r.Header.Get("Upgrade") == "websocket" {
			log.Printf("Proxying WebSocket: %s %s", r.Method, r.URL.Path)

			// Determine backend WebSocket URL from BACKEND_URL
			backendWSURL := strings.Replace(BACKEND_URL, "http://", "ws://", 1)
			backendWSURL = strings.Replace(backendWSURL, "https://", "wss://", 1)
			backendWSURL = backendWSURL + "/ws"

			// Upgrade client connection
			clientConn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Printf("WebSocket upgrade error: %v", err)
				return
			}
			defer clientConn.Close()

			// Prepare headers for backend connection (forward auth headers)
			// Create a clean header set to avoid duplicates
			backendHeaders := make(http.Header)

			// Forward authentication headers
			if cookie := r.Header.Get("Cookie"); cookie != "" {
				backendHeaders.Set("Cookie", cookie)
			}

			// DO NOT forward Sec-WebSocket-* headers - the Dialer adds them automatically
			// Forwarding them causes "duplicate header not allowed" errors
			// Only forward Origin header if needed
			for key, values := range r.Header {
				keyLower := strings.ToLower(key)
				if keyLower == "origin" {
					// Forward origin if present
					if len(values) > 0 {
						backendHeaders.Set(key, values[0])
					}
				}
			}

			// Connect to backend WebSocket
			backendConn, _, err := websocket.DefaultDialer.Dial(backendWSURL, backendHeaders)
			if err != nil {
				log.Printf("Backend WebSocket connection error: %v", err)
				clientConn.WriteMessage(websocket.CloseMessage, []byte("Backend connection failed"))
				return
			}
			defer backendConn.Close()

			// Proxy messages in both directions
			errChan := make(chan error, 2)

			// Client -> Backend
			go func() {
				for {
					messageType, message, err := clientConn.ReadMessage()
					if err != nil {
						errChan <- err
						return
					}
					if err := backendConn.WriteMessage(messageType, message); err != nil {
						errChan <- err
						return
					}
				}
			}()

			// Backend -> Client
			go func() {
				for {
					messageType, message, err := backendConn.ReadMessage()
					if err != nil {
						errChan <- err
						return
					}
					if err := clientConn.WriteMessage(messageType, message); err != nil {
						errChan <- err
						return
					}
				}
			}()

			// Wait for an error from either direction
			<-errChan
			log.Printf("WebSocket proxy connection closed")
			return
		}

		// Regular HTTP request
		proxy.ServeHTTP(w, r)
	}

	// Main handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Proxy API, WebSocket, and uploads requests to backend
		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/ws" || strings.HasPrefix(r.URL.Path, "/uploads/") {
			log.Printf("Proxying: %s %s", r.Method, r.URL.Path)

			// Special handling for WebSocket
			if r.URL.Path == "/ws" {
				wsProxy(w, r)
			} else {
				proxy.ServeHTTP(w, r)
			}
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

		// Add no-cache headers for CSS files to prevent caching issues during development
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}

		fs.ServeHTTP(w, r)
	})

	log.Printf("🚀 Frontend server running at http://localhost%s", PORT)
	log.Printf("🔄 Proxying /api/* and /ws to %s", BACKEND_URL)
	log.Printf("📁 Serving static files from: ./client")
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
