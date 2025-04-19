package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
	"github.com/akhmadst1/tugas-akhir-backend/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func CreateDocument(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse uploaded file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		// Use the filename as the title (without extension)
		fileName := file.Filename
		title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		if title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
			return
		}

		// Save the file to the "docs" directory
		filePath := fmt.Sprintf("docs/%s", file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Create the document entry in DB
		doc := models.Document{
			Title:  title,
			Source: filePath,
		}

		if err := repositories.CreateDocument(db, &doc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
			return
		}

		c.JSON(http.StatusCreated, doc)
	}
}

func ViewDocument(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
			return
		}

		doc, err := repositories.GetDocumentByID(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			return
		}

		c.Header("Content-Type", "application/pdf")
		c.File(doc.Source)
	}
}

func GetAllDocuments(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		docs, err := repositories.GetAllDocuments(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
			return
		}

		c.JSON(http.StatusOK, docs)
	}
}

func GetDocumentByID(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
			return
		}

		doc, err := repositories.GetDocumentByID(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			return
		}

		c.JSON(http.StatusOK, doc)
	}
}

func UpdateDocument(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var doc models.Document
		if err := c.ShouldBindJSON(&doc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
			return
		}
		doc.ID = id

		if err := repositories.UpdateDocument(db, &doc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Document updated"})
	}
}

func DeleteDocument(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
			return
		}

		if err := repositories.DeleteDocument(db, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Document deleted"})
	}
}
