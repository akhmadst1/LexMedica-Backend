package handlers

import (
	"net/http"

	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
	"github.com/akhmadst1/tugas-akhir-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// HandleLLMAnalysis forwards the text to the Disharmony LLM API
func HandleDisharmonyAnalysis(c *gin.Context) {
	var req models.DisharmonyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	response, err := services.GetDisharmonyAnalysis(req.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze text"})
		return
	}

	c.JSON(http.StatusOK, response)
}
