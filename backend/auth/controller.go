package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// SetupAuthRoutes configures the routes used for OAuth flow
func (c *AuthController) SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.GET("/login", c.Login)
		auth.GET("/callback", c.Callback)
	}
}

// Login initiates the OAuth2 flow
func (c *AuthController) Login(ctx *gin.Context) {
	state := "random-state" // In production, use a secure random state
	url := c.authService.GetAuthURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback handles the OAuth2 callback
func (c *AuthController) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	userInfo, err := c.authService.HandleCallback(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Here you would typically:
	// 1. Create or update the user in your database
	// 2. Generate a JWT or session token
	// 3. Set the token in a cookie or return it to the frontend

	ctx.JSON(http.StatusOK, gin.H{
		"user": userInfo,
	})
}
