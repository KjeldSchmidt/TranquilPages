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
	url, err := c.authService.GetAuthURL()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate redirect url for OAuth flow"})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback handles the OAuth2 callback
func (c *AuthController) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	state := ctx.Query("state")
	if state == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "State parameter is required"})
		return
	}

	userInfo, err := c.authService.HandleCallback(code, state)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := GenerateToken(userInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token after login"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":  userInfo,
		"token": token,
	})
}
