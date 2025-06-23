package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	configs "ai-mockup/config"
	"ai-mockup/repository"
	"ai-mockup/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  configs.AppConfig.GoogleAuth.RedirectURL,
	ClientID:     configs.AppConfig.GoogleAuth.ClientID,
	ClientSecret: configs.AppConfig.GoogleAuth.ClientSecret,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func GetGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  configs.AppConfig.GoogleAuth.RedirectURL,
		ClientID:     configs.AppConfig.GoogleAuth.ClientID,
		ClientSecret: configs.AppConfig.GoogleAuth.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func OautHandler(c *gin.Context) {
	oauthconfig := GetGoogleOauthConfig()
	url := oauthconfig.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func OauthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	oauthconfig := GetGoogleOauthConfig()
	token, err := oauthconfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to exchange token",
		})
		return
	}

	client := oauthconfig.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to get user info",
		})
		return
	}
	defer userInfoResp.Body.Close()

	var userInfo struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to parse user data",
		})
		return
	}

	user, err := repository.GetUserByEmail(userInfo.Email)
	if err != nil {
		// User does not exist, create new
		user, err = repository.CreateUser(userInfo.Name, userInfo.Email, userInfo.Picture, "google")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to create user",
			})
			return
		}
	}

	tokenStr, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to generate token",
		})
		return
	}

	cookieStr := "__did_ddy__dick__l_e__er__prod=" + tokenStr +
		"; Path=/; Domain=.localhost:3000; Max-Age=86400; HttpOnly; Secure; SameSite=Lax"
	c.Header("Set-Cookie", cookieStr)
	c.Redirect(http.StatusPermanentRedirect, configs.AppConfig.DashboardURL)
}

func LogoutHandler(c *gin.Context) {
	logoutCookie := "__did_ddy__dick__l_e__er__prod=; Path=/; Domain=.localhost:3000; Max-Age=-1; HttpOnly; Secure; SameSite=Lax"
	c.Header("Set-Cookie", logoutCookie)
	c.Redirect(http.StatusPermanentRedirect, configs.AppConfig.DashboardURL)
}

func SessionHandler(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}
	user, err := repository.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"name":     user.Name,
			"avatar":   user.Avatar,
			"provider": user.Provider,
		},
	})
}
