package repositories

import (
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// CreateDocument inserts a new document
func CreateDocument(db *sqlx.DB, doc *models.Document) error {
	query := `INSERT INTO documents (title, source) VALUES ($1, $2) RETURNING id`
	return db.QueryRow(query, doc.Title, doc.Source).Scan(&doc.ID)
}

// GetAllDocuments retrieves all documents
func GetAllDocuments(db *sqlx.DB) ([]models.Document, error) {
	var docs []models.Document
	query := `SELECT * FROM documents ORDER BY id DESC`
	err := db.Select(&docs, query)
	return docs, err
}

// GetDocumentByID retrieves a single document by ID
func GetDocumentByID(db *sqlx.DB, id int) (*models.Document, error) {
	var doc models.Document
	query := `SELECT * FROM documents WHERE id = $1`
	err := db.Get(&doc, query, id)
	return &doc, err
}

// UpdateDocument updates an existing document
func UpdateDocument(db *sqlx.DB, doc *models.Document) error {
	query := `UPDATE documents SET title = $1, source = $2 WHERE id = $3`
	_, err := db.Exec(query, doc.Title, doc.Source, doc.ID)
	return err
}

// DeleteDocument deletes a document by ID
func DeleteDocument(db *sqlx.DB, id int) error {
	query := `DELETE FROM documents WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
