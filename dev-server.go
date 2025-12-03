// dev-server.go - Simple development server for testing the SPA
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := "3000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Serve static files from client/ directory
	fs := http.FileServer(http.Dir("./client"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("./client", r.URL.Path)

		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// File doesn't exist, serve index.html (SPA routing)
			http.ServeFile(w, r, "./client/index.html")
			return
		}

		// File exists, serve it
		fs.ServeHTTP(w, r)
	})

	fmt.Printf("🚀 Dev server running at http://localhost:%s\n", port)
	fmt.Println("📁 Serving from: ./client")
	fmt.Println("Press Ctrl+C to stop")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
