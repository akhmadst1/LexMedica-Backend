package models

type ChatSession struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	StartedAt string `json:"started_at"`
}

type ChatMessage struct {
	ID                   int                   `json:"id"`
	SessionID            int                   `json:"session_id"`
	Sender               string                `json:"sender"`
	Message              string                `json:"message"`
	CreatedAt            string                `json:"created_at"`
	ProcessingTimeMs     int                   `json:"processing_time_ms"`
	DisharmonyAnalysis   []DisharmonyAnalysis  `json:"disharmony_analysis"` // now it's a slice!
	ChatMessageDocuments []ChatMessageDocument `json:"chat_message_documents"`
}

type DisharmonyAnalysis struct {
	ID               int    `json:"id"`
	MessageID        int    `json:"message_id"`
	Result           bool   `json:"result"`
	Analysis         string `json:"analysis"`
	ProcessingTimeMs int    `json:"processing_time_ms"`
}
