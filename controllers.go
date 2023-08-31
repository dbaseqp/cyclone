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
		username := strings.Split(formdata.SessionKey, ":")[0]
		//password := GeneratePassword(12)

		_, err := CloneOnDemand(formdata, username)
		// username = fmt.Sprintf("%s_%s", username, pg)

		if err != nil {
			c.JSON(401, gin.H{"message": err.Error()})
			// c.Abort()
			return
		}
		// c.JSON(200, gin.H{"message": gin.H{"username": username, "password": password}})
		c.JSON(200, gin.H{"message": "success"})
	}
}

func TemplateGuestViewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formdata models.AuthForm

		if c.ShouldBindJSON(&formdata) != nil {
			c.JSON(401, gin.H{"message": "Missing fields"})
			// c.Abort()
			return
		}

		creds := strings.Split(formdata.SessionKey, ":")

		templates, err := TemplateGuestView(creds[0], creds[1])

		if err != nil {
			templates, err = TemplateGuestView(tomlConf.VCenterUsername, tomlConf.VCenterPassword)
			fmt.Println(err)
			if err != nil {
				c.JSON(500, gin.H{"message": "Problem getting templates"})
				// c.Abort()
				return
			}
		}
		c.JSON(200, gin.H{"message": templates})
	}
}

func PodViewHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formdata models.AuthForm

		if c.ShouldBindJSON(&formdata) != nil {
			c.JSON(401, gin.H{"message": "Missing fields"})
			// c.Abort()
			return
		}

		username, err := ValidateJWT(formdata.SessionKey)

		if err != nil {
			c.JSON(401, gin.H{"message": "Unauthorized"})
			return
		}

		userPods, err := ViewPods(username)

		if err != nil {
			c.JSON(500, gin.H{"message": "Problem getting pods"})
			// c.Abort()
			return
		}

		c.JSON(200, gin.H{"message": userPods})
	}
}

func DeletePodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formdata models.DeletePodForm

		if c.ShouldBindJSON(&formdata) != nil {
			c.JSON(406, gin.H{"message": "Missing fields"})
			// c.Abort()
			return
		}
		// formdata = models.InvokeCloneOnDemandForm{
		// 	Template: c.PostForm("template"),
		// }
		// username := "Bruharmy"
		username := strings.Split(formdata.SessionKey, ":")[0]

		err := DeletePod(formdata, username)
		// username = fmt.Sprintf("%s_%s", username, pg)

		if err != nil {
			c.JSON(401, gin.H{"message": "Problem deleting pod"})
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
		_, err := ValidateJWT(formdata.SessionKey)
		if err != nil {
			c.JSON(401, gin.H{"message": "Unauthorized"})
			return
		}
		c.JSON(200, gin.H{"message": "success"})
	}
}
