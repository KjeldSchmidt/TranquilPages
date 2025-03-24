package auth

import (
	"log"
	"net/http"
	"os"

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
		auth.POST("/logout", c.Logout)
	}

	api := router.Group("/api")
	api.GET("/user/me", AuthMiddleware(c.authService), c.GetCurrentUser)
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

	// Set token in HTTP-only cookie
	ctx.SetCookie("token", token, 3600, "/", "", true, true)

	frontendURL := os.Getenv("FRONTEND_URL")
	frontendURL, exists := os.LookupEnv("FRONTEND_URL")
	if !exists {
		log.Fatal("Failed to get FRONTEND_URL from environment")
	}

	ctx.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// Logout handles user logout by blacklisting their token
func (c *AuthController) Logout(ctx *gin.Context) {
	token, err := getTokenFromRequest(ctx)
	if err != nil {
		switch err.(type) {
		case *TokenNotFoundError, *InvalidAuthHeaderError:
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	if err := c.authService.Logout(token); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	// Clear the cookie
	ctx.SetCookie("token", "", -1, "/", "", true, true)
	ctx.Status(http.StatusNoContent)
}

// GetCurrentUser returns the current user's information
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userClaims, ok := claims.(*Claims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":      userClaims.UserID,
		"email":   userClaims.Email,
		"name":    userClaims.Name,
		"picture": userClaims.Picture,
	})
}
