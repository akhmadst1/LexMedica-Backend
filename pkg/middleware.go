package pkg

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for Bearer token in Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid authorization"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		if _, err := ValidateJWT(token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid authorization"})
			c.Abort()
			return
		}

		// Token is valid, continue to handler
		c.Next()
	}
}
