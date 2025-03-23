package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// getTokenFromRequest extracts the authentication token from either cookie or Authorization header
func getTokenFromRequest(c *gin.Context) (string, error) {
	// First try to get token from cookie
	token, err := c.Cookie("token")
	if err == nil {
		return token, nil
	}

	// If no cookie, try Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", &TokenNotFoundError{}
	}

	// Check if it's a Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", &InvalidAuthHeaderError{}
	}

	return parts[1], nil
}

// AuthMiddleware checks for a valid session token
func AuthMiddleware(authService *AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c)
		if err != nil {
			switch err.(type) {
			case *TokenNotFoundError, *InvalidAuthHeaderError:
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			c.Abort()
			return
		}

		claims, err := authService.ValidateAuthenticationToken(token)
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
