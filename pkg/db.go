package pkg

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectDB() *sqlx.DB {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Create tables if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS users (
	    id SERIAL PRIMARY KEY,
	    email VARCHAR(255) UNIQUE NOT NULL,
	    password TEXT NOT NULL,
	    verified BOOLEAN DEFAULT FALSE,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS chat_sessions (
	    id SERIAL PRIMARY KEY,
	    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	    title VARCHAR(255) NOT NULL,
	    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS documents (
	    id SERIAL PRIMARY KEY,
	    title VARCHAR(255) NOT NULL,
	    source TEXT NOT NULL,
	    content BYTEA NOT NULL
	);

	CREATE TABLE IF NOT EXISTS chat_messages (
    	id SERIAL PRIMARY KEY,
    	session_id INTEGER NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    	sender VARCHAR(10) CHECK (sender IN ('bot', 'user')) NOT NULL,
    	message TEXT NOT NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS chat_message_documents (
    	message_id INTEGER NOT NULL REFERENCES chat_messages(id) ON DELETE CASCADE,
    	document_id INTEGER NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    	PRIMARY KEY (message_id, document_id)
	);`

	db.MustExec(schema)

	return db
}
