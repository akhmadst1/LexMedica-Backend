package repositories

import (
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

// CreateUser stores a new user
func CreateUser(db *sqlx.DB, user models.User) error {
	_, err := db.Exec("INSERT INTO users (email, password, verified) VALUES ($1, $2, $3)", user.Email, user.Password, user.Verified)
	return err
}

// Verify email account
func VerifyUser(db *sqlx.DB, email string) error {
	query := "UPDATE users SET verified = TRUE WHERE email = $1"
	_, err := db.Exec(query, email)
	return err
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(db *sqlx.DB, email string) (*models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT email, password, verified FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err // Return nil instead of &user when no user is found
	}
	return &user, nil
}
