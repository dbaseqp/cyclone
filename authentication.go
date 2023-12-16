package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// authRequired provides authentication middleware for ensuring that a user is logged in.

// getUUID returns a randomly generated UUID
func getUUID() string {
	return uuid.NewString()
}

// initCookies use gin-contrib/sessions{/cookie} to initalize a cookie store.
// It generates a random secret for the cookie store -- not ideal for continuity or invalidating previous cookies, but it's secure and it works
func initCookies(router *gin.Engine) {
	router.Use(sessions.Sessions("kamino", cookie.NewStore([]byte("kamino")))) // change to secret
}

func getUser(c *gin.Context) string {
	userID := sessions.Default(c).Get("Username")
	if userID != nil {
		return userID.(string)
	}
	return ""
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("Username")
	if id == nil {
		c.JSON(http.StatusOK, gin.H{"message": "No session."})
		return
	}
	session.Delete("Username")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out!"})
}

func validateAgainstSSO(c *gin.Context) {
	token, err := c.Request.Cookie("auth_token")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url := fmt.Sprintf("https://%s/api/users/auth/%s", os.Getenv("WSSO_FQDN"), token.Value)

	resp, err := http.Get(url)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Next()
}

func tokenAuth(c *gin.Context) {
	token := c.Query("token")
	fmt.Println(token)
	fmt.Println(c.Request.URL)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	url := fmt.Sprintf("https://%s/api/users/auth/%s", os.Getenv("WSSO_FQDN"), token)

	resp, err := http.Get(url)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}

	var response map[string]string

	err = json.Unmarshal([]byte(w.body.String()), &response)

	c.SetCookie("auth_token", token, 86400, "/", "kamino.dev.gfed", false, false)

	session := sessions.Default(c)
	session.Set("Username", response["Username"])
	session.Set("Groups", response["Groups"])
	session.Set("Admin", response["Admin"])

	c.Redirect(http.StatusFound, "https://kamino.dev.gfed/dashboard")
}
