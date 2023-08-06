package main
import (
	// "log"
	// "net/http"
	// "strconv"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	// jwt "github.com/dgrijalva/jwt-go"

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
		// username := "Bruharmy"
		username := strings.Split(formdata.SessionKey,":")[0]
		password := GeneratePassword(12)

		_, err := CloneOnDemand(formdata, username, password)
		// username = fmt.Sprintf("%s_%s", username, pg)
		
		if err != nil {
			c.JSON(401, gin.H{"message": "Problem cloning pod"})
			// c.Abort()
			return
		}
		// c.JSON(200, gin.H{"message": gin.H{"username": username, "password": password}})
		c.JSON(200, gin.H{"message": "success"})
	}
}

func UserRegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formdata models.LoginForm
	
		if c.ShouldBindJSON(&formdata) != nil {
			c.JSON(406, gin.H{"message": "Missing fields"})
			// c.Abort()
			return
		}

		err := RegisterUser(formdata.Username, formdata.Password)

		if err != nil {
			fmt.Println(err)
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "success"})
	}
}

func UserLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formdata models.LoginForm
	
		if c.ShouldBindJSON(&formdata) != nil {
			c.JSON(406, gin.H{"message": "Missing fields"})
			// c.Abort()
			return
		}
		// user := &models.User{
		// 	Username: "test",
		// 	Role: "guest",
		// }
		// token, err := GenerateToken(user.Username)
		// if err != nil {
		// 	c.JSON(401, gin.H{"message": "Invalid session"})
		// }

		if ValidateVsphereCredentials(formdata.Username, formdata.Password) != nil {
			c.JSON(401, gin.H{"message": "Invalid credentials"})
			return
		}
		token := fmt.Sprintf("%s:%s", formdata.Username, formdata.Password)
		c.JSON(200, gin.H{"message": token})
	}
}

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formdata models.AuthForm

		if c.ShouldBindJSON(&formdata) != nil {
			c.JSON(401, gin.H{"message": "Missing fields"})
			// c.Abort()
			return
		}

		// token, err := GenerateToken(user.Username)
		// if err != nil {
		// 	c.JSON(401, gin.H{"message": "Invalid session"})
		// }
		if ValidateJWT(formdata.SessionKey) != nil {
			c.JSON(401, gin.H{"message": "Unauthorized"})
			return
		}
		c.JSON(200, gin.H{"message": "success"})
	}
}
