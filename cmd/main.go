package main

import (
	"log"
	"time"

	"github.com/akhmadst1/tugas-akhir-backend/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Router
	r := gin.Default()

	// Serve static files for test cases few shot examples
	r.Static("/data", "./data")

	// ** Add CORS Middleware **
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache preflight request for 12 hours
	}))

	// ** API Routes **
	r.POST("/api/chat/analyze", handlers.DisharmonyAnalysis)

	// Start Server
	r.Run(":8080")
}
