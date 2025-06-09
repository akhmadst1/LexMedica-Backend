package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"
)

var supabaseJWKS jwk.Set
var jwtAudience = "authenticated" // Default audience in Supabase
var jwtIssuer string              // Will be your Supabase URL

func InitAuth(supabaseURL string) {
	jwtIssuer = supabaseURL + "/auth/v1"
	var err error
	supabaseJWKS, err = jwk.Fetch(context.Background(), jwtIssuer+"/.well-known/jwks.json")
	if err != nil {
		panic("Failed to fetch JWKS: " + err.Error())
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			kid := token.Header["kid"].(string)
			key, ok := supabaseJWKS.LookupKeyID(kid)
			if !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			var pubKey interface{}
			if err := key.Raw(&pubKey); err != nil {
				return nil, err
			}
			return pubKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			return
		}

		// Attach user ID to context
		sub, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing sub in token"})
			return
		}
		c.Set("user_id", sub)
		c.Next()
	}
}
