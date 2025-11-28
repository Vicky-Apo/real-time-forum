// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"frontend-service/config"
	"frontend-service/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	app, err := routes.NewApp(cfg)
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	mux := routes.SetupRoutes(app)

	fmt.Printf("Frontend server starting on port %s\n", cfg.FrontendBaseURL)
	fmt.Printf("Backend API URL: %s\n", cfg.APIBaseURL)

	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
