package files

import (
	"github.com/gin-gonic/gin"
)

func CreateRouterGroup(r *gin.RouterGroup) {
	files := r.Group("/files")

	{
		files.POST("", HandleFilesUpload)
		files.GET("", HandleGetAllFiles)
		files.GET("/:fileName", HandleFilesDownload)
		files.HEAD("/:fileName", HandleFilesSize)
	}
}
