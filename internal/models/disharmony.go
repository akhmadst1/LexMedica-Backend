package models

// DisharmonyRequest represents a regulation text for analysis
type DisharmonyRequest struct {
	Text string `json:"text" binding:"required"`
}

// DisharmonyResponse represents the disharmony analysis result
type DisharmonyResponse struct {
	Result string `json:"result"`
}
