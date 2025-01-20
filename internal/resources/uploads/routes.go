package uploads

import (
	"github.com/gin-gonic/gin"
)

func CreateRouterGroup(r *gin.RouterGroup) {
	uploads := r.Group("/uploads")

	{
		uploads.POST("", HandleFilesUpload)
		uploads.GET("", HandleGetAllFiles)
	}
}
