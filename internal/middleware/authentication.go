package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"knights-vow/pkg/jwt"
)

func AuthenticateUser(c *gin.Context) {
	if c.Request.URL.Path == "/api/v1/users/login" || c.Request.URL.Path == "/api/v1/users/sign-up" {
		c.Next()
		return
	}

	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(401, gin.H{
			"error": "Authorization header missing",
		})

		c.Abort()
		return
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")

	if !jwt.Verify(tokenString) {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})

		c.Abort()
		return
	}

	c.Next()
}
