package files

import (
	"knights-vow/pkg/path"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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

func HandleGetAllFiles(c *gin.Context) {
	files, err := GetAllFiles()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting files",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "files retrieved",
		"count":   len(files),
		"files":   files,
	})
}

func HandleFilesDownload(c *gin.Context) {
	fileName := c.Param("fileName")

	filePath, err := path.CreatePathFromRoot("data/uploads/" + fileName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting file path",
			"error":   err,
		})
		return
	}

	file, err := os.Open(filePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error opening file",
			"error":   err,
		})
		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting file info",
			"error":   err,
		})
		return
	}

	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(int(fileInfo.Size())))

	http.ServeContent(c.Writer, c.Request, fileName, fileInfo.ModTime(), file)
}

func HandleFilesSize(c *gin.Context) {
	fileName := c.Param("fileName")

	filePath, err := path.CreatePathFromRoot("data/uploads/" + fileName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting file path",
			"error":   err,
		})
		return
	}

	file, err := os.Open(filePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error opening file",
			"error":   err,
		})
		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting file info",
			"error":   err,
		})
		return
	}

	c.Writer.Header().Set("Content-Length", strconv.Itoa(int(fileInfo.Size())))
	c.Status(http.StatusOK)
}
