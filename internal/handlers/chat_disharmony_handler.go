package handlers

import (
	"net/http"

	"github.com/akhmadst1/tugas-akhir-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

func CreateChatDisharmony(c *gin.Context) {
	var chatDisharmonyRequest struct {
		MessageID int    `json:"message_id"`
		Result    bool   `json:"result"`
		Analysis  string `json:"analysis"`
	}

	if err := c.BindJSON(&chatDisharmonyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	analysis, err := repositories.CreateChatDisharmony(chatDisharmonyRequest.MessageID, chatDisharmonyRequest.Result, chatDisharmonyRequest.Analysis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert disharmony analysis to chat"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}
