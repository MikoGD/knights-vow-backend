package files

import (
	"encoding/json"
	"fmt"
	"io"
	"knights-vow/internal/resources/users"
	"knights-vow/pkg/path"
	"knights-vow/pkg/sockets"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	fileName := c.Param("fileName")

	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing fileName query parameter",
		})
		return
	}

	outputPath, err := path.CreatePathFromRoot("data/uploads/" + fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong",
			"error":   err,
		})
		return
	}

	output, err := os.Create(outputPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong",
			"error":   err,
		})
		return
	}

	defer output.Close()

	_, err = io.Copy(output, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong",
			"error":   err,
		})
		return
	}

	c.JSON(204, gin.H{
		"message": "file uploaded successfully",
	})
}

func HandleFileDownload(c *gin.Context, filesService FilesService) {
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
	// UPDATE: Use createErrorResponse
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse file ID",
			"error":   err,
		})
		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "Cannot read file")
		return
	}

	// Set headers
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", string(rune(fileInfo.Size())))

	// Stream the file content
	c.File(file.Name())
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
