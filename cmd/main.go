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
		AllowOrigins:     []string{"http://localhost:3000"}, // Update with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache preflight request for 12 hours
	}))

	// ** API Routes **
	// Auth
	r.POST("/auth/register", handlers.RegisterUser(db))
	r.GET("/auth/verify", handlers.VerifyUser(db))
	r.POST("/auth/login", handlers.LoginUser(db))
	r.POST("/auth/refresh", handlers.RefreshToken(db))

	// Chat History
	r.POST("/history/sessions", handlers.CreateChatSession(db))
	r.GET("/history/sessions", handlers.GetChatSessions(db))
	r.DELETE("/history/sessions/:id", handlers.DeleteChatSession(db))

	r.POST("/history/messages", handlers.AddMessage(db))
	r.GET("/history/messages/:session_id", handlers.GetMessages(db))
	r.DELETE("/history/messages/:id", handlers.DeleteMessage(db))

	// QnA and Analysis
	r.POST("/qna", handlers.HandleQnARequest)
	r.POST("/analyze", handlers.HandleDisharmonyAnalysis)

	// Start Server
	log.Println("Web App Backend running on port 8080...")
	r.Run(":8080")
}
