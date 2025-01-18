package users

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	myjwt "knights-vow/pkg/jwt"
)

func handleCreateUser(c *gin.Context) {
	userPayload := &UserPayload{}

	var err error

	err = c.Bind(userPayload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request payload",
			"error":   err,
		})
		return
	}

	user, err := GetUserByUsername(userPayload.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	userID, err := SaveUser(userPayload.Username, string(hashedPassword))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"token":   myjwt.CreateJWT(userID),
		"URI":     "/users/" + strconv.Itoa(userID),
	})
}

func HandleUserLogin(c *gin.Context) {
	userPayload := &UserPayload{}

	var err error

	err = c.Bind(userPayload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request payload",
			"error":   err,
		})
		return
	}

	user, err := GetUserByUsername(userPayload.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid login or password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPayload.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid login or password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user logged in",
		"token":   myjwt.CreateJWT(user.ID),
	})
}

func CheckUserAuthStatus(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user ID",
		})
		return
	}

	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Authorization header missing",
		})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token := myjwt.ParseJWT(tokenString)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["client_id"] != float64(userID) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid client_id in JWT",
			})
			return
		}
	}

	if token.Valid {
		c.JSON(http.StatusOK, gin.H{
			"message":         "user is authenticated",
			"isAuthenticated": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "user is authenticated",
		"isAuthenticated": false,
	})
}
