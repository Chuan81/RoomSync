package main

import (
	"fmt"
	"log"
	"roomsync/api"
	"roomsync/config"
	"roomsync/repository"
)

func main() {
	// Initialize Configuration
	config.InitConfig()

	// Initialize Database
	repository.InitDB()

	// Setup Gin Router
	r := api.SetupRouter()

	// Start Server
	port := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("RoomSync server starting on %s...", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
