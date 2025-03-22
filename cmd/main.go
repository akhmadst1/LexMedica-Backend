package main

import (
	"log"

	"github.com/akhmadst1/tugas-akhir-backend/internal/handlers"
	"github.com/akhmadst1/tugas-akhir-backend/pkg"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Load env variables
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to Database
	db := pkg.ConnectDB()
	defer db.Close()

	// Initialize Router
	r := gin.Default()

	// ** API Routes
	// Auth
	r.POST("/register", handlers.RegisterUser(db))
	r.POST("/login", handlers.LoginUser(db))
	r.POST("/refresh", handlers.RefreshToken(db))

	// QnA and Analysis
	r.POST("/qna", handlers.HandleQnARequest)
	r.POST("/analyze", handlers.HandleDisharmonyAnalysis)

	// Start Server
	log.Println("Web App Backend running on port 8080...")
	r.Run(":8080")
}
