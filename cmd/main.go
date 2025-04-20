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

	// Serve static files for the documents
	r.Static("/docs", "./docs")

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
	r.POST("/auth/resend_email_verification", handlers.ResendEmailVerification(db))
	r.GET("/auth/verify_email/:token", handlers.VerifyEmailUser(db))
	r.POST("/auth/login", handlers.LoginUser(db))
	r.POST("/auth/refresh_token", handlers.RefreshToken(db))

	// QnA and Analysis
	r.POST("/chat/qna", handlers.HandleQnARequest)
	r.POST("/chat/analyze", handlers.HandleDisharmonyAnalysis)

	// History routes require authentication
	history := r.Group("/history", pkg.AuthMiddleware())
	{
		history.POST("/session", handlers.CreateChatSession(db))
		history.GET("/session/:user_id", handlers.GetChatSessions(db))
		history.DELETE("/session/:session_id", handlers.DeleteChatSession(db))

		history.POST("/message", handlers.AddChatMessage(db))
		history.GET("/message/:session_id", handlers.GetChatMessages(db))
		history.DELETE("/message/:message_id", handlers.DeleteMessage(db))
	}

	// Documents routes require authentication
	docs := r.Group("/document", pkg.AuthMiddleware())
	{
		docs.POST("", handlers.CreateDocument(db))
		docs.GET("", handlers.GetAllDocuments(db))
		docs.GET("/:id", handlers.GetDocumentByID(db))
		docs.PUT("/:id", handlers.UpdateDocument(db))
		docs.DELETE("/:id", handlers.DeleteDocument(db))
		docs.GET("/view/:id", handlers.ViewDocument(db))
	}

	// Start Server
	log.Println("Web App Backend running on port 8080...")
	r.Run(":8080")
}
