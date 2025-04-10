package handlers

import (
	"net/http"

	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
	"github.com/akhmadst1/tugas-akhir-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// HandleQnARequest forwards the question to the QnA microservice
func HandleQnARequest(c *gin.Context) {
	var req models.QnARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Setup SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Writer.WriteHeader(http.StatusOK)
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}

	// Start streaming
	if err := services.StreamQnAAnswer(req.Question, c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Streaming error: %v", err)
	}
}
