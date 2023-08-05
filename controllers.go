package main
import (
	// "log"
	// "net/http"
	// "strconv"
	"fmt"

	"github.com/gin-gonic/gin"
	jwt "github.com/dgrijalva/jwt-go"

	"bruharmy/models"
)

func PingGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	}
}

func TemplateGuestViewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		templates, err := TemplateGuestView()
		
		if err != nil {
			c.JSON(500, gin.H{"message": "Problem getting templates"})
			// c.Abort()
			return
		}
		c.JSON(200, gin.H{"message": templates})
	}
}

func CloneOnDemandHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formdata models.InvokeCloneOnDemandForm
		
		if c.ShouldBindJSON(&formdata) != nil {
			c.JSON(406, gin.H{"message": "Missing fields"})
			// c.Abort()
			return
		}
		// formdata = models.InvokeCloneOnDemandForm{
		// 	Template: c.PostForm("template"),
		// }
		username := "Bruharmy"
		password := GeneratePassword(12)

		pg, err := CloneOnDemand(formdata, username, password)
		
		if err != nil {
			c.JSON(400, gin.H{"message": "Problem cloning pod"})
			// c.Abort()
			return
		}
		c.JSON(200, gin.H{"message": gin.H{"username": fmt.Sprintf("%s_%s", username, pg), "password": password}})
	}

	func UserRegisterHandler() gin.HandlerFunc {
		return func(c *gin.Context) {
			
		}
	}

	func UserLoginHandler() gin.HandlerFunc {
		return func(c *gin.Context) {
			user := &models.User{
				Username: "test",
				Role: "guest"
			}
			token, err := GenerateToken(user.Username)
			if err != nil {
				c.JSON(401, gin.H{"message": "Invalid session"})
			}
			c.JSON(200, gin.H{"message": token})
		}
	}
}