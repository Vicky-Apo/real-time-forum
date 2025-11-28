package main

import (
	"fmt"
	"log"
	"net/http"

	"platform.zone01.gr/git/gpapadopoulos/forum/config"
	"platform.zone01.gr/git/gpapadopoulos/forum/database"
	"platform.zone01.gr/git/gpapadopoulos/forum/internal/routes"
)

func main() {
	// Load configuration
	err := config.LoadConfig()
	if err != nil {
		log.Panic(err)
	}

	// Initialize the database
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Setup API routes
	apiRoutes := routes.SetupRoutes(db)

	// Create server using config values
	serverAddr := fmt.Sprintf("%s:%s", config.Config.ServerHost, config.Config.ServerPort)

	fmt.Printf("Starting API server on %s\n", serverAddr)

	log.Fatal(http.ListenAndServe(serverAddr, apiRoutes))
}
