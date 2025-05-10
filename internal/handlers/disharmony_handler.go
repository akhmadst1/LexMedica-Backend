package handlers

import (
	"net/http"

	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
	"github.com/akhmadst1/tugas-akhir-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// Disharmony Regulations Analyzer
func DisharmonyAnalysis(c *gin.Context) {
	var req models.DisharmonyAnalysisRequest
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

	// OpenAI Model
	if err := services.StreamOpenAIDisharmonyAnalysis(req.Regulations, c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Streaming error: %v", err)
	}
}
