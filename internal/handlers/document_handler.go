package handlers

import (
	"net/http"

	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
	"github.com/akhmadst1/tugas-akhir-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

func CreateChatDocuments(c *gin.Context) {
	var createRequests []struct {
		MessageID  int    `json:"message_id"`
		DocumentID int    `json:"document_id"`
		Clause     string `json:"clause"`
		Snippet    string `json:"snippet"`
		PageNumber int    `json:"page_number"`
	}

	if err := c.BindJSON(&createRequests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var docs []models.ChatMessageDocument
	for _, req := range createRequests {
		docs = append(docs, models.ChatMessageDocument{
			MessageID:  req.MessageID,
			DocumentID: req.DocumentID,
			Clause:     req.Clause,
			Snippet:    req.Snippet,
			PageNumber: req.PageNumber,
		})
	}

	insertedDocs, err := repositories.CreateChatDocuments(docs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create documents"})
		return
	}

	c.JSON(http.StatusOK, insertedDocs)
}

func GetLinkDocumentByTypeNumberYear(c *gin.Context) {
	documentType := c.Param("type") // e.g., api/document/:type/:number/:year
	number := c.Param("number")
	year := c.Param("year")
	if documentType == "" || number == "" || year == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing type or number or year"})
		return
	}

	sessions, err := repositories.GetLinkDocumentByTypeNumberYear(documentType, number, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document by type, number, and year"})
		return
	}
	c.JSON(http.StatusOK, sessions)
}
