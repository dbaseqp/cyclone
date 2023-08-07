package main

import (
	"github.com/gin-gonic/gin"
)

func PublicRoutes(g *gin.RouterGroup) {
	// account/authentication endpoints
	g.POST("/account/login", UserLoginHandler())
	g.POST("/account/register", UserRegisterHandler())
	g.POST("/account/auth", AuthHandler())

	// vsphere pod endpoints
	g.GET("/pods/templates", TemplateGuestViewHandler())
	g.POST("/pods/clone", CloneOnDemandHandler())
	g.POST("/pods/view", PodViewHandler())
	g.POST("/pods/delete", DeletePodHandler())

}

func PrivateRoutes(g *gin.RouterGroup) {
	
}
