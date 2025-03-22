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

// GetUserByEmail retrieves a user by email
func GetUserByEmail(db *sqlx.DB, email string) (*models.User, error) {
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	return &user, err
}

// UpdateRefreshToken updates the refresh token of a user
func UpdateRefreshToken(db *sqlx.DB, email, refreshToken string) error {
	_, err := db.Exec("UPDATE users SET refresh_token=$1 WHERE email=$2", refreshToken, email)
	return err
}
