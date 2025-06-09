package main

import (
	"time"

	"github.com/akhmadst1/tugas-akhir-backend/config"
	"github.com/akhmadst1/tugas-akhir-backend/internal/handlers"
	"github.com/akhmadst1/tugas-akhir-backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Supabase client
	config.Init()

	// Initialize Router
	r := gin.Default()

	// Serve static files for test cases few shot examples
	r.Static("/data", "./data")

	// ** Add CORS Middleware **
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://lex-medica-frontend.vercel.app"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ** API Routes **
	chat := r.Group("/api/chat")
	{
		session := chat.Group("/session")
		session.Use(middleware.RequireAuth())
		{
			session.POST("", handlers.CreateChatSession)
			session.GET("/:user_id", handlers.GetChatSessionsByUserID)
			session.DELETE("/:id", handlers.DeleteChatSession)
		}

		message := chat.Group("/message")
		message.Use(middleware.RequireAuth())
		{
			message.POST("", handlers.CreateChatMessage)
			message.GET("/:session_id", handlers.GetChatMessagesBySessionID)
		}

		disharmony := chat.Group("/disharmony")
		disharmony.Use(middleware.RequireAuth())
		{
			disharmony.POST("", handlers.CreateChatDisharmony)
		}

		document := chat.Group("/document")
		document.Use(middleware.RequireAuth())
		{
			document.POST("", handlers.CreateChatDocuments)
		}

		chat.POST("/qna", handlers.QNA)
		chat.POST("/analyze", handlers.DisharmonyAnalysis)
	}

	document := r.Group("/api/document")
	{
		document.GET("/:type/:number/:year", handlers.GetLinkDocumentByTypeNumberYear)
	}

	// Start Server
	r.Run(":8080")
}
