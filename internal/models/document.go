package models

type Document struct {
	ID     int    `json:"id" db:"id"`
	Title  string `json:"title" db:"title"`
	Source string `json:"source" db:"source"`
}
