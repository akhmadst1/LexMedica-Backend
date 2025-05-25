package main

import (
	"log"
	"time"

	"github.com/akhmadst1/tugas-akhir-backend/config"
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

	// Initialize Supabase client
	config.Init()

	// Initialize Router
	r := gin.Default()

	// Serve static files for test cases few shot examples
	r.Static("/data", "./data")

	// ** Add CORS Middleware **
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache preflight request for 12 hours
	}))

	// ** API Routes **
	chat := r.Group("/api/chat")
	{
		session := chat.Group("/session")
		{
			session.POST("", handlers.CreateChatSession)
			session.GET("/:user_id", handlers.GetChatSessionsByUserID)
			session.DELETE("/:id", handlers.DeleteChatSession)
		}

		message := chat.Group("/message")
		{
			message.POST("", handlers.CreateChatMessage)
			message.GET("/:session_id", handlers.GetChatMessagesBySessionID)
		}

		disharmony := chat.Group("/disharmony")
		{
			disharmony.POST("", handlers.CreateChatDisharmony)
		}

		// chat.POST("/qna", handlers.QNA)
		chat.POST("/analyze", handlers.DisharmonyAnalysis)
	}

	// Start Server
	r.Run(":8080")
}
