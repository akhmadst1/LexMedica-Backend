package models

type ChatSession struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	StartedAt string `json:"started_at"`
}

type ChatMessage struct {
	ID         int             `json:"id,omitempty"`
	SessionID  int             `json:"session_id,omitempty"`
	Sender     string          `json:"sender"`
	Message    string          `json:"message"`
	CreatedAt  string          `json:"created_at,omitempty"`
	Disharmony *ChatDisharmony `json:"disharmony,omitempty"`
	Documents  []Document      `json:"documents,omitempty"`
}

type ChatDisharmony struct {
	ID        int    `json:"id,omitempty"`
	MessageID int    `json:"message_id,omitempty"`
	Result    bool   `json:"result"`
	Analysis  string `json:"analysis"`
	CreatedAt string `json:"created_at,omitempty"`
}
