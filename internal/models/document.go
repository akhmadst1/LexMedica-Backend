package models

type ChatMessageDocument struct {
	MessageID    int          `json:"message_id"`
	Clause       string       `json:"clause"`
	DocumentID   int          `json:"document_id"`
	Snippet      string       `json:"snippet"`
	LinkDocument LinkDocument `json:"link_documents"`
}

type LinkDocument struct {
	ID     int    `json:"id,omitempty"`
	About  string `json:"about,omitempty"`
	Type   string `json:"type,omitempty"`
	Number int    `json:"number,omitempty"`
	Year   int    `json:"year,omitempty"`
	Status string `json:"status,omitempty"`
	URL    string `json:"url,omitempty"`
}
