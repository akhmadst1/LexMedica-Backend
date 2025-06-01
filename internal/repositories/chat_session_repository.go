package repositories

import (
	"fmt"
	"sort"

	"github.com/akhmadst1/tugas-akhir-backend/config"
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
)

func CreateChatSession(userID string, title string) (models.ChatSession, error) {
	var insertedSessions []models.ChatSession

	err := config.Supabase.DB.
		From("chat_sessions").
		Insert(map[string]interface{}{
			"user_id": userID,
			"title":   title,
		}).
		Execute(&insertedSessions)

	if err != nil || len(insertedSessions) == 0 {
		return models.ChatSession{}, fmt.Errorf("insert failed or no rows returned")
	}

	return insertedSessions[0], nil
}

func GetChatSessionsByUserID(userID string) ([]models.ChatSession, error) {
	var sessions []models.ChatSession
	err := config.Supabase.DB.
		From("chat_sessions").
		Select("*").
		Eq("user_id", userID).
		Execute(&sessions)

	// Reverse the order of sessions to have the most recent first
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].ID > sessions[j].ID
	})

	return sessions, err
}

func DeleteChatSession(sessionID string) error {
	return config.Supabase.DB.
		From("chat_sessions").
		Delete().
		Eq("id", sessionID).
		Execute(nil)
}
