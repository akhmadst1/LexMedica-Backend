package repositories

import (
	"fmt"

	"github.com/akhmadst1/tugas-akhir-backend/config"
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
)

func CreateChatDocuments(docs []models.ChatMessageDocument) ([]models.ChatMessageDocument, error) {
	if len(docs) == 0 {
		return nil, fmt.Errorf("no documents to insert")
	}

	// Convert each struct to map[string]interface{} for Supabase insert
	records := make([]map[string]interface{}, len(docs))
	for i, doc := range docs {
		records[i] = map[string]interface{}{
			"message_id":  doc.MessageID,
			"document_id": doc.DocumentID,
			"clause":      doc.Clause,
			"snippet":     doc.Snippet,
			"page_number": doc.PageNumber,
		}
	}

	var inserted []models.ChatMessageDocument
	err := config.Supabase.DB.
		From("chat_message_documents").
		Insert(records).
		Execute(&inserted)

	if err != nil || len(inserted) == 0 {
		return nil, fmt.Errorf("failed to insert chat message documents: %w", err)
	}

	return inserted, nil
}

func GetLinkDocumentByTypeNumberYear(documentType string, number string, year string) ([]models.LinkDocument, error) {
	var linkDocument []models.LinkDocument

	err := config.Supabase.DB.
		From("link_documents").
		Select(`*`).
		Eq("type", documentType).
		Eq("number", number).
		Eq("year", year).
		Execute(&linkDocument)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch link document: %w", err)
	}

	return linkDocument, nil
}
