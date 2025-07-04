package files

import (
	"fmt"
	"knights-vow/pkg/sockets"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	PathToTempDirectory = "data/temp"
	PathToUploads       = "data/uploads"
	ChunkSize           = 1024 * 1024
)

func HandleGetFiles(c *gin.Context, filesService FilesService) {
	fileName := c.Query("fileName")

	files, err := filesService.GetFiles(fileName)

	if err != nil {
		status, errResponse := createErrorResponse(err)
		c.JSON(status, errResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "files retrieved",
		"count":   len(files),
		"files":   files,
	})
}

func HandleFileUpload(c *gin.Context, filesService FilesService) {
	ws, err := sockets.Upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error upgrading connection",
			"error":   err,
		})
	}

	filesService.UploadFile(ws)
}

func HandleFileUpload2(c *gin.Context, filesService FilesService) {
	fileName := c.Query("fileName")
	ownerIDQueryParam := c.Query("ownerID")

	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing fileName query parameter",
		})
		return
	}

	if ownerIDQueryParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing fileName query parameter",
		})
		return
	}

	ownerID, err := strconv.Atoi(ownerIDQueryParam)
	if err != nil {
		status, errorResponse := createErrorResponse(err)
		c.JSON(status, errorResponse)
		return
	}

	err = filesService.SaveFileFromUpload(fileName, ownerID, c.Request.Body)
	if err != nil {
		status, errorResponse := createErrorResponse(err)
		c.JSON(status, errorResponse)
		return
	}

	c.JSON(204, gin.H{
		"message": "file uploaded successfully",
	})
}

func HandleFileDownload(c *gin.Context, filesService FilesService) {
	log.Println("Handle file download")
	fileIDParam := c.Param("fileID")
	fileID, err := strconv.Atoi(fileIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse file ID",
			"error":   err,
		})
		return
	}

	file, err := filesService.GetFileForDownload(fileID)
	if err != nil {
		status, errorResponse := createErrorResponse(err)
		c.JSON(status, errorResponse)
		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		status, errorResponse := createErrorResponseWithCustomDefaultMesage(err, "Failed to read file")
		c.JSON(status, errorResponse)
		return
	}

	log.Printf("Sending file %s of length %d\n", fileInfo.Name(), fileInfo.Size())

	c.DataFromReader(http.StatusOK, fileInfo.Size(), "application/octet-stream", file, map[string]string{
		"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name()),
	})
}

func HandleDeleteFile(c *gin.Context, filesService FilesService) {
	fileIDParam := c.Param("fileID")

	fileID, err := strconv.Atoi(fileIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid file ID",
		})
		return
	}

	if err = filesService.DeleteFile(fileID); err != nil {
		status, errorResponse := createErrorResponse(err)
		c.JSON(status, errorResponse)
	}

	c.Status(http.StatusNoContent)
}
