package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/akhmadst1/tugas-akhir-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func QNA(c *gin.Context) {
	var req services.QnaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	apiKey := os.Getenv("QNA_API_KEY")

	data, err := services.QNAService(req, apiKey)
	if err != nil {
		log.Printf("Error calling QnA service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch from QnA service"})
		return
	}

	// Use c.Data to return raw JSON if already serialized
	c.Data(http.StatusOK, "application/json", data)
}
