package users

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"

	"knights-vow/pkg/jwt"
)

func handleCreateUser(c *gin.Context) {
	userPayload := &UserPayload{}

	var err error

	err = c.Bind(userPayload)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request payload",
			"error":   err,
		})
		return
	}

	user, err := GetUserByUsername(userPayload.Username)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	if user != nil {
		c.JSON(409, gin.H{
			"error": gin.H{
				"username": "username already exists",
			},
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPayload.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	userID, err := SaveUser(userPayload.Username, string(hashedPassword))

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "user created",
		"token":   jwt.CreateJWT(),
		"URI":     "/users/" + strconv.Itoa(userID),
	})
}

func HandleUserLogin(c *gin.Context) {
	userPayload := &UserPayload{}

	var err error

	err = c.Bind(userPayload)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid request payload",
			"error":   err,
		})
		return
	}

	user, err := GetUserByUsername(userPayload.Username)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	if user == nil {
		c.JSON(401, gin.H{
			"error": "invalid login or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPayload.Password))

	if err != nil {
		c.JSON(401, gin.H{
			"error": "invalid login or password",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "user logged in",
		"token":   jwt.CreateJWT(),
	})
}
