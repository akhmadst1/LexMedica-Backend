package models

type ChatMessageDocument struct {
	MessageID    int          `json:"message_id"`
	Clause       string       `json:"clause"`
	DocumentID   int          `json:"document_id"`
	Snippet      string       `json:"snippet"`
	LinkDocument LinkDocument `json:"link_documents"`
}

type LinkDocument struct {
	ID     int    `json:"id"`
	About  string `json:"about"`
	Type   string `json:"type"`
	Number int    `json:"number"`
	Year   int    `json:"year"`
	Status string `json:"status"`
	URL    string `json:"url"`
}
