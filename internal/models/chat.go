package models

import "time"

type ChatSession struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	StartedAt time.Time `json:"started_at" db:"started_at"`
}

type ChatMessage struct {
	ID        int       `json:"id" db:"id"`
	SessionID int       `json:"session_id" db:"session_id"`
	Sender    string    `json:"sender" db:"sender"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
