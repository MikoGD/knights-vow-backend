package files

import (
	"database/sql"
	"knights-vow/internal/resources/users"

	"github.com/gin-gonic/gin"
)

func CreateRouterGroup(r *gin.RouterGroup, db *sql.DB) {
	files := r.Group("/files")
	usersRepository := users.CreateNewUserRepository(db)
	filesRepository := CreateNewFilesRepository(db, usersRepository)
	filesService := createNewFilesService(filesRepository)

	{
		files.GET("", func(c *gin.Context) {
			HandleGetFiles(c, filesService)
		})

		files.GET("/upload", func(c *gin.Context) {
			HandleFileUpload(c, filesService)
		})

		files.GET("/:fileID", func(c *gin.Context) {
			HandleFileDownload(c, filesService)
		})

		files.DELETE("/:fileID", func(c *gin.Context) {
			HandleDeleteFile(c, filesService)
		})
	}
}
