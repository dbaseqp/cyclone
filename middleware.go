package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := ValidateToken(c)
		if err != nil {
			c.String(401, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}