package files

import (
	"database/sql"
	"knights-vow/internal/resources/users"

	"github.com/gin-gonic/gin"
)

func CreateRouterGroup(r *gin.RouterGroup, db *sql.DB) {
	files := r.Group("/files")
	usersRepository := users.CreateNewUserRepository(db)

	{
		files.GET("", func(c *gin.Context) {
			HandleGetFiles(c, db)
		})

		files.GET("/upload", func(c *gin.Context) {
			HandleFileUpload(c, db, usersRepository)
		})

		files.GET("/:fileID", func(c *gin.Context) {
			HandleFileDownload(c, db)
		})

		files.DELETE("/:fileID", func(c *gin.Context) {
			HandleDeleteFile(c, db)
		})
	}
}
