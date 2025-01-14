package users

import (
	"github.com/gin-gonic/gin"
)

func CreateRouterGroup(r *gin.RouterGroup) {
	users := r.Group("/users")

	{
		users.POST("/sign-up", handleCreateUser)
		users.POST("/login", HandleUserLogin)
	}
}
