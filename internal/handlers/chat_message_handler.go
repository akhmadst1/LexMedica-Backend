package handlers

import (
	"net/http"

	"github.com/akhmadst1/tugas-akhir-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

func CreateChatMessage(c *gin.Context) {
	var chatMessageRequest struct {
		SessionID int    `json:"session_id"`
		Sender    string `json:"sender"`
		Message   string `json:"message"`
	}

	if err := c.BindJSON(&chatMessageRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	message, err := repositories.CreateChatMessage(chatMessageRequest.SessionID, chatMessageRequest.Sender, chatMessageRequest.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert message"})
		return
	}

	c.JSON(http.StatusOK, message)
}

func GetChatMessagesBySessionID(c *gin.Context) {
	sessionID := c.Param("session_id")
	messages, err := repositories.GetChatMessagesBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
