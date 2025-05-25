package models

type ChatDocument struct {
	MessageID  int          `json:"message_id,omitempty"`
	DocumentID int          `json:"document_id,omitempty"`
	Clause     string       `json:"clause"`
	Snippet    string       `json:"snippet"`
	Source     LinkDocument `json:"source"`
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
