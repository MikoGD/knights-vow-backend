package files

import (
	"github.com/gin-gonic/gin"
)

func CreateRouterGroup(r *gin.RouterGroup) {
	files := r.Group("/files")

	{
		files.GET("", HandleGetAllFiles)
		files.GET("/upload", HandleFileUpload)
		files.GET("/:fileID", HandleFileDownload)
	}
}
