package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid session token
func AuthMiddleware(authService *AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := authService.ValidateAuthenticationToken(parts[1])
		if err != nil {
			if _, ok := err.(*TokenRevokedError); ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			}
			c.Abort()
			return
		}

		// Set claims in context for downstream handlers
		c.Set("claims", claims)
		c.Next()
	}
}
