package models

type ChatSession struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	StartedAt string `json:"started_at"`
}

type ChatMessage struct {
	ID                 int                 `json:"id,omitempty"`
	SessionID          int                 `json:"session_id,omitempty"`
	Sender             string              `json:"sender"`
	Message            string              `json:"message"`
	CreatedAt          string              `json:"created_at,omitempty"`
	DisharmonyAnalysis *DisharmonyAnalysis `json:"disharmony_analysis,omitempty"`
	Documents          []ChatDocument      `json:"documents,omitempty"`
}

type DisharmonyAnalysis struct {
	ID        int    `json:"id,omitempty"`
	MessageID int    `json:"message_id,omitempty"`
	Result    bool   `json:"result"`
	Analysis  string `json:"analysis"`
}
