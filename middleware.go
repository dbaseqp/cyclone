package main

import (
	"github.com/gin-gonic/gin"
)

func JwtAuthRequired(c *gin.Context) {
	err := ValidateToken(c)
	if err != nil {
		c.String(401, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
}