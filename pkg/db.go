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

	// Create table if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS users (
	    id SERIAL PRIMARY KEY,
	    email VARCHAR(255) UNIQUE NOT NULL,
	    password TEXT NOT NULL,
	    verified BOOLEAN DEFAULT FALSE,
	    refresh_token TEXT,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	db.MustExec(schema)

	return db
}
