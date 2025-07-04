package files

import (
	"database/sql"
	"knights-vow/internal/resources/users"
	"log"

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

		files.POST("", func(c *gin.Context) {
			HandleFileUpload2(c, filesService)
		})

		files.GET("/:fileID", func(c *gin.Context) {
			log.Println("Route /:fileID")
			HandleFileDownload(c, filesService)
		})

		files.DELETE("/:fileID", func(c *gin.Context) {
			HandleDeleteFile(c, filesService)
		})
	}
}
