package handlers

import (
	"net/http"

	"github.com/akhmadst1/tugas-akhir-backend/pkg"

	"github.com/akhmadst1/tugas-akhir-backend/internal/models"

	"github.com/akhmadst1/tugas-akhir-backend/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.User
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		req.Password, _ = pkg.HashPassword(req.Password)
		err := repositories.CreateUser(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
			return
		}

		token, _ := pkg.GenerateJWT(req.Email)
		pkg.SendVerificationEmail(req.Email, token)

		c.JSON(http.StatusOK, gin.H{"message": "Check your email for verification"})
	}
}

func VerifyUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from query parameters
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
			return
		}

		// Validate token
		email, err := pkg.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Update user as verified in the database
		err = repositories.VerifyUser(db, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Account successfully verified"})
	}
}

func LoginUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Retrieve user from database
		user, err := repositories.GetUserByEmail(db, req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Check if user is verified
		if !user.Verified {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account not verified"})
			return
		}

		// Validate password
		if !pkg.CheckPassword(user.Password, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate tokens
		jwtToken, _ := pkg.GenerateJWT(user.Email)
		refreshToken, _ := pkg.GenerateRefreshToken(user.Email)

		// Update refresh token in the database
		repositories.UpdateRefreshToken(db, user.Email, refreshToken)

		// Return tokens
		c.JSON(http.StatusOK, gin.H{
			"jwt":          jwtToken,
			"refreshToken": refreshToken,
		})
	}
}

func RefreshToken(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refreshToken"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Validate refresh token
		email, err := pkg.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		// Generate new JWT token
		newJwtToken, _ := pkg.GenerateJWT(email)

		c.JSON(http.StatusOK, gin.H{
			"jwt": newJwtToken,
		})
	}
}
