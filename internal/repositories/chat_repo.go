package repositories

import (
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// CreateChatSession creates a new chat session and returns the created session
func CreateChatSession(db *sqlx.DB, session *models.ChatSession) error {
	query := `INSERT INTO chat_sessions (user_id, title) VALUES ($1, $2) RETURNING id, started_at`
	return db.QueryRow(query, session.UserID, session.Title).Scan(&session.ID, &session.StartedAt)
}

// GetChatSessions retrieves all chat sessions for a user
func GetChatSessions(db *sqlx.DB, userID int) ([]models.ChatSession, error) {
	var sessions []models.ChatSession
	query := `SELECT * FROM chat_sessions WHERE user_id = $1`
	err := db.Select(&sessions, query, userID)
	return sessions, err
}

// DeleteChatSession deletes a chat session by ID
func DeleteChatSession(db *sqlx.DB, sessionID int) error {
	query := `DELETE FROM chat_sessions WHERE id = $1`
	_, err := db.Exec(query, sessionID)
	return err
}

// AddChatMessage adds a new message to a chat session
func AddChatMessage(db *sqlx.DB, message *models.ChatMessage) error {
	query := `INSERT INTO chat_messages (session_id, sender, message) VALUES ($1, $2, $3) RETURNING id, created_at`
	return db.QueryRow(query, message.SessionID, message.Sender, message.Message).Scan(&message.ID, &message.CreatedAt)
}

// GetChatMessages retrieves all messages from a chat session
func GetChatMessages(db *sqlx.DB, sessionID int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	query := `SELECT * FROM chat_messages WHERE session_id = $1 ORDER BY created_at ASC`
	err := db.Select(&messages, query, sessionID)
	return messages, err
}

// DeleteMessage deletes a chat message by ID
func DeleteMessage(db *sqlx.DB, messageID int) error {
	query := `DELETE FROM chat_messages WHERE id = $1`
	_, err := db.Exec(query, messageID)
	return err
}
