package repositories

import (
	"fmt"
	"strings"

	"github.com/akhmadst1/tugas-akhir-backend/config"
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
)

func CreateChatDocument(messageID int, documentID int, clause string, snippet string) (models.ChatDocument, error) {
	var insertedChatDocument []models.ChatDocument

	err := config.Supabase.DB.
		From("chat_message_documents").
		Insert(map[string]interface{}{
			"message_id":  messageID,
			"document_id": documentID,
			"clause":      clause,
			"snippet":     snippet,
		}).
		Execute(&insertedChatDocument)

	if err != nil || len(insertedChatDocument) == 0 {
		return models.ChatDocument{}, fmt.Errorf("failed to insert chat message document: %w", err)
	}

	return insertedChatDocument[0], nil
}

func GetLinkDocumentByTypeNumberYear(documentType string, number string, year string) ([]models.LinkDocument, error) {
	var linkDocument []models.LinkDocument

	err := config.Supabase.DB.
		From("link_documents").
		Select(`*`).
		Eq("type", strings.ToUpper(documentType)).
		Eq("number", number).
		Eq("year", year).
		Execute(&linkDocument)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch link document: %w", err)
	}

	return linkDocument, nil
}
