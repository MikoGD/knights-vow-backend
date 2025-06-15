package users

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	myjwt "knights-vow/pkg/jwt"
)

func handleCreateUser(c *gin.Context, userService UserService) {
	userPayload := &userDTO{}

	var err error

	err = c.Bind(userPayload)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request payload",
			"error":   err,
		})
		return
	}

	userID, err := userService.CreateUser(userPayload)
	if err != nil {
		status, errorResponse := createErrorResponse(err)
		c.JSON(status, errorResponse)
		return
	}

	c.JSON(http.StatusCreated, &createUserResponse{
		message: "Successfully created user",
		token:   myjwt.CreateJWT(userID),
	})
}

func HandleUserLogin(c *gin.Context, userService UserService) {
	userPayload := &userDTO{}

	err := c.Bind(userPayload)

	token, err := userService.AuthenticateUser(userPayload)
	if err != nil {
		status, errorResponse := createErrorResponse(err)
		c.JSON(status, errorResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged in user",
		"token":   token,
	})
}

func CheckUserAuthStatus(c *gin.Context, userService UserService) {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user ID",
		})
		return
	}

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Authorization header missing",
		})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	isAuthenticated, err := userService.CheckUserAuthStatus(tokenString, userID)
	if err != nil {
		status, errorResponse := createErrorResponse(err)
		c.JSON(status, errorResponse)
		return
	}

	if isAuthenticated {
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusUnauthorized)
}
