package models

// QnARequest represents a question from the frontend
type QnARequest struct {
	Question string `json:"question" binding:"required"`
}

// QnAResponse represents the answer from the RAG service
type QnAResponse struct {
	Answer string `json:"answer"`
}
