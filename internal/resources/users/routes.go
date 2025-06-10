package users

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func CreateRouterGroup(r *gin.RouterGroup, db *sql.DB) {
	users := r.Group("/users")
	repository := CreateNewUserRepository(db)
	userService := createNewUserService(repository)

	{
		users.POST("/sign-up", func(c *gin.Context) {
			handleCreateUser(c, userService)
		})
		users.POST("/login", func(c *gin.Context) {
			HandleUserLogin(c, userService)
		})
		users.GET("/:id/auth-status", func(c *gin.Context) {
			CheckUserAuthStatus(c, userService)
		})
	}
}
