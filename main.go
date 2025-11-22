package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"InfluenceIQ/config"
	"InfluenceIQ/middleware"
	"InfluenceIQ/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("  Warning: .env file not found, using system environment")
	}

	// Connect to DB
	config.ConnectDB()
	defer config.CloseDB()

	// Initialize router
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Create /api group for all routes
	api := router.Group("/api")
	routes.RegisterAuthRoutes(api)

	// Start server
	log.Println(" Server running on http://localhost:8080")
	router.Run(":8080")
}
