package handlers

import (
	"net/http"

	"github.com/akhmadst1/tugas-akhir-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

func CreateChatSession(c *gin.Context) {
	var CreateChatSessionRequest struct {
		UserID string `json:"user_id"`
		Title  string `json:"title"`
	}

	if err := c.BindJSON(&CreateChatSessionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	session, err := repositories.CreateChatSession(CreateChatSessionRequest.UserID, CreateChatSessionRequest.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func GetChatSessionsByUserID(c *gin.Context) {
	userID := c.Param("user_id") // e.g., api/chat/sessions/:user_id
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id"})
		return
	}

	sessions, err := repositories.GetChatSessionsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat sessions"})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func DeleteChatSession(c *gin.Context) {
	sessionID := c.Param("id") // DELETE api/chat/sessions/:id

	if err := repositories.DeleteChatSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat session deleted successfully"})
}
