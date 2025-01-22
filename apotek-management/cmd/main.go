package main

import (
	"apotek-management/config"
	"apotek-management/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Initialize a new Gin router
	router := gin.Default()

	// Set up the routes using the SetupRoutes function
	routes.SetupRoutes(router)

	return router
}

func main() {
	// Connect to the database
	config.ConnectDB()

	// Start the server with the routes
	r := setupRouter()
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
