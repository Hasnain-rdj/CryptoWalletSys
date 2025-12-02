package main

import (
	"backend/config"
	"backend/middleware"
	"backend/routes"
	"backend/services"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize MongoDB
	if err := config.InitMongoDB(); err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer config.DisconnectMongoDB()

	// Create database indexes for optimal performance
	if err := config.CreateIndexes(); err != nil {
		log.Printf("Warning: Failed to create some indexes: %v", err)
	}

	// Initialize blockchain with genesis block
	services.InitBlockchain()

	// Start Zakat scheduler
	go services.StartZakatScheduler()

	// Setup Gin router
	r := gin.Default()

	// Global middleware
	r.Use(middleware.SanitizeMiddleware())

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(r)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
