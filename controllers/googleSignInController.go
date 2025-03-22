package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var clientId string = os.Getenv("GOOGLE_CLIENT_ID")
var ClientSecret string = os.Getenv("GOOGLE_CLIENT_SECRET")
var redirectUrl string = os.Getenv("REDIRECT_URL")

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: ClientSecret,
		RedirectURL:  redirectUrl,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
)

func HandleGoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("google-auth-state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "google-auth-state" {
		fmt.Printf("invalid oauth state, expected 'random-state', got '%s'\n", state)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("Code exchange failed with '%s'\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		fmt.Printf("Failed getting user info: %s\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer response.Body.Close()

	fmt.Printf("Response status: %s\n", response.Status)

	c.String(http.StatusOK, "Google Sign in successfull")
}
