package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"knights-vow/pkg/jwt"
)

func AuthenticateUser(c *gin.Context) {
	var token string

	if len(c.Request.Header["Upgrade"]) == 1 {
		token = c.Query("token")
	} else if c.Request.URL.Path != "/api/v1/users/login" && c.Request.URL.Path != "/api/v1/users/sign-up" {
		token = c.GetHeader("Authorization")
	} else {
		c.Next()
		return
	}

	if token == "" {
		log.Printf("Failed to get token when accessing path: %s\n", c.Request.URL.Path)
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
