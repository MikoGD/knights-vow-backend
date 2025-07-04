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

	response := createUserResponse{
		Message: "Successfully created user",
		Token:   myjwt.CreateJWT(userID),
	}

	c.JSON(http.StatusCreated, response)
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

	response := loginUserResponse{
		Message: "Successfully logged in user",
		Token:   token,
	}

	c.JSON(http.StatusOK, response)
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
		c.JSON(http.StatusOK, authStatusResponse{
			IsAuthenticated: true,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, authStatusResponse{
		IsAuthenticated: false,
	})
}
