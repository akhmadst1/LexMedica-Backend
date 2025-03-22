package main

import (
	"log"

	"github.com/akhmadst1/tugas-akhir-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Define API routes
	r.POST("/qna", handlers.HandleQnARequest)
	r.POST("/analyze", handlers.HandleDisharmonyAnalysis)

	// Start server
	log.Println("Web App Backend running on port 8080...")
	r.Run(":8080")
}
