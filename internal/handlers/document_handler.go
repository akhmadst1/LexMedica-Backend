package handlers

import (
	"net/http"

	"github.com/akhmadst1/tugas-akhir-backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

func CreateChatDocument(c *gin.Context) {
	var CreateChatDocumentRequest struct {
		MessageID  int    `json:"message_id"`
		DocumentID int    `json:"document_id"`
		Clause     string `json:"clause"`
		Snippet    string `json:"snippet"`
	}

	if err := c.BindJSON(&CreateChatDocumentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	document, err := repositories.CreateChatDocument(CreateChatDocumentRequest.MessageID,
		CreateChatDocumentRequest.DocumentID, CreateChatDocumentRequest.Clause, CreateChatDocumentRequest.Snippet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}

	c.JSON(http.StatusOK, document)
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
