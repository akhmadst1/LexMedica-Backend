package main

import (
	"log"
	"time"

	"github.com/akhmadst1/tugas-akhir-backend/internal/handlers"
	"github.com/akhmadst1/tugas-akhir-backend/pkg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	// ** Add CORS Middleware **
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Update with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache preflight request for 12 hours
	}))

	// ** API Routes **
	// Auth
	r.POST("/register", handlers.RegisterUser(db))
	r.GET("/verify", handlers.VerifyUser(db))
	r.POST("/login", handlers.LoginUser(db))
	r.POST("/refresh", handlers.RefreshToken(db))

	// QnA and Analysis
	r.POST("/qna", handlers.HandleQnARequest)
	r.POST("/analyze", handlers.HandleDisharmonyAnalysis)

	// Start Server
	log.Println("Web App Backend running on port 8080...")
	r.Run(":8080")
}
