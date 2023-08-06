package main

import (
	"github.com/gin-gonic/gin"
)

func PublicRoutes(g *gin.RouterGroup) {
	// g.GET("/login", controllers.LoginGetHandler())
	// g.POST("/login", controllers.LoginPostHandler())
	// g.GET("/", controllers.IndexGetHandler())

	g.GET("/templates/guest", TemplateGuestViewHandler())
	g.POST("/login", UserLoginHandler())
	g.POST("/register", UserRegisterHandler())
	g.POST("/api/auth", AuthHandler())
	g.POST("/clone/ondemand", CloneOnDemandHandler())
	g.POST("/pod/delete", DeletePodHandler())
}

func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/ping", PingGetHandler())
}
