package uploads

import (
	"knights-vow/pkg/path"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HandleFilesUpload(c *gin.Context) {
	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error parsing form",
			"error":   err,
		})
		return
	}

	files := form.File["files"]
	username := form.Value["username"]

	if len(username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username not provided",
		})
		return
	}

	uploadsFilePath, err := path.CreatePathFromRoot("data/uploads")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating uploads directory",
			"error":   err,
		})
		return
	}

	fileNames := make([]string, len(files))

	for i, file := range files {
		err = c.SaveUploadedFile(file, filepath.Join(uploadsFilePath, file.Filename))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error saving file",
				"error":   err,
			})
			return
		}

		fileNames[i] = file.Filename
	}

	filesUploaded, err := SaveFiles(fileNames, username[0])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error saving file to database",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "file(s) uploaded",
		"filesUploaded": filesUploaded,
	})
}
