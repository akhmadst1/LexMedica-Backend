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
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error":  "Invalid input",
			})
			return
		}

		// Check if user already exists
		existingUser, err := repositories.GetUserByEmail(db, req.Email)
		if err == nil && existingUser != nil {
			c.JSON(http.StatusConflict, gin.H{
				"status": http.StatusConflict,
				"error":  "Email is already registered",
			})
			return
		}

		// Hash password
		hashedPassword, err := pkg.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Failed to secure password",
			})
			return
		}
		req.Password = hashedPassword

		// Attempt to create user
		err = repositories.CreateUser(db, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Could not create user",
			})
			return
		}

		// Generate token and send verification email
		token, err := pkg.GenerateJWT(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error":  "Could not generate verification token",
			})
			return
		}
		pkg.SendVerificationEmail(req.Email, token)

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Check registered email for verification",
		})
	}
}

func ResendEmailVerification(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		user, err := repositories.GetUserByEmail(db, req.Email)
		if err != nil || user.Verified {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request or already verified"})
			return
		}

		// Generate a new token and resend
		token, err := pkg.GenerateJWT(user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		pkg.SendVerificationEmail(user.Email, token)

		c.JSON(http.StatusOK, gin.H{"message": "Verification email resent"})
	}
}

func VerifyEmailUser(db *sqlx.DB) gin.HandlerFunc {
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
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error":  "Invalid input",
			})
			return
		}

		// Retrieve user from database
		user, err := repositories.GetUserByEmail(db, req.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"error":  "User not registered",
			})
			return
		}

		// Check if user is verified
		if !user.Verified {
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"error":  "Account not verified",
			})
			return
		}

		// Validate password
		if !pkg.CheckPassword(user.Password, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  "Incorrect email or password",
			})
			return
		}

		// Generate tokens
		jwtToken, _ := pkg.GenerateJWT(user.Email)
		refreshToken, _ := pkg.GenerateRefreshToken(user.Email)

		// Send tokens to frontend
		c.JSON(http.StatusOK, gin.H{
			"status":       http.StatusOK,
			"message":      "Login successful",
			"id":           user.ID,
			"email":        user.Email,
			"token":        jwtToken,
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

		// Validate refresh token (no need to check DB)
		email, err := pkg.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		// Generate new JWT token
		newJwtToken, _ := pkg.GenerateJWT(email)

		c.JSON(http.StatusOK, gin.H{
			"token": newJwtToken,
		})
	}
}
