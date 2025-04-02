package handlers

import (
	"net/http"
	"strconv"

	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
	"github.com/akhmadst1/tugas-akhir-backend/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// CreateChatSession creates a new chat session
func CreateChatSession(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var session models.ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := repositories.CreateChatSession(db, &session); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			return
		}

		c.JSON(http.StatusCreated, session)
	}
}

// GetChatSessions retrieves all chat sessions for a user
func GetChatSessions(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Query("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		sessions, err := repositories.GetChatSessions(db, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sessions"})
			return
		}

		c.JSON(http.StatusOK, sessions)
	}
}

// DeleteChatSession deletes a chat session
func DeleteChatSession(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		if err := repositories.DeleteChatSession(db, sessionID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Chat session deleted"})
	}
}

// AddMessage adds a new message to a chat session
func AddMessage(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var message models.ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := repositories.AddMessage(db, &message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add message"})
			return
		}

		c.JSON(http.StatusCreated, message)
	}
}

// GetMessages retrieves all messages from a chat session
func GetMessages(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := strconv.Atoi(c.Param("session_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		messages, err := repositories.GetMessages(db, sessionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		c.JSON(http.StatusOK, messages)
	}
}

// DeleteMessage deletes a chat message
func DeleteMessage(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		messageID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
			return
		}

		if err := repositories.DeleteMessage(db, messageID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Chat message deleted"})
	}
}
