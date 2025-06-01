package repositories

import (
	"fmt"
	"sort"

	"github.com/akhmadst1/tugas-akhir-backend/config"
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
)

func CreateChatMessage(sessionID int, sender string, message string) (models.ChatMessage, error) {
	var insertedMessages []models.ChatMessage

	err := config.Supabase.DB.
		From("chat_messages").
		Insert(map[string]interface{}{
			"session_id": sessionID,
			"sender":     sender,
			"message":    message,
		}).
		Execute(&insertedMessages)

	if err != nil || len(insertedMessages) == 0 {
		return models.ChatMessage{}, fmt.Errorf("failed to insert chat message: %w", err)
	}

	return insertedMessages[0], nil
}

func GetChatMessagesBySessionID(sessionID string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage

	err := config.Supabase.DB.
		From("chat_messages").
		Select("id,session_id,sender,message,created_at,disharmony_analysis(*),chat_message_documents(message_id,document_id,clause,snippet,link_documents(*))").
		Eq("session_id", sessionID).
		Execute(&messages)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch chat messages: %w", err)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].ID < messages[j].ID
	})

	return messages, nil
}
