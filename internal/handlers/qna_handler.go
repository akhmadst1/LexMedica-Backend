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

	response, err := services.GetQnAAnswer(req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch answer"})
		return
	}

	c.JSON(http.StatusOK, response)
}
