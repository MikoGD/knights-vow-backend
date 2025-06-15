package files

import (
	"database/sql"
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

func HandleGetFiles(c *gin.Context, db *sql.DB) {
	fileName := c.Query("fileName")

	var files []File
	var err error

	if fileName == "" {
		files, err = GetAllFiles(db)
	} else {
		files, err = GetFilesByName(fileName, db)
	}

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

func HandleFileUpload(c *gin.Context, db *sql.DB, usersRepository users.UserRepository) {
	ws, err := sockets.Upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error upgrading connection",
			"error":   err,
		})
	}

	_, payload, err := ws.ReadMessage()

	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error reading message")
		return
	}

	var initMessage FileUploadInitMessage
	err = json.Unmarshal(payload, &initMessage)
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error unmarshalling message")
		return
	}

	tempDir, err := CreateTempDir(initMessage.UserID, initMessage.FileName)
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error creating temp directory")
		return
	}

	chunksCount := 0
	for i := 1; i <= initMessage.TotalChunks; i++ {
		_, payload, err = ws.ReadMessage()
		if err != nil {
			sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error reading message")
			return
		}

		err := SaveChunk(tempDir, i, payload)

		if err != nil {
			sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error saving chunk")
		}

		percentageUploaded := math.Round((float64(i) / float64(initMessage.TotalChunks)) * 100)

		err = ws.WriteJSON(gin.H{
			"message":          "Chunk saved",
			"chunkNumber":      i,
			"uploadPercentage": int(percentageUploaded),
		})

		if err != nil {
			sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error sending save chunk message")
		}

		chunksCount++
	}

	finalFilePath, err := path.CreatePathFromRoot("data/uploads/" + initMessage.FileName)
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error creating final file path")
		return
	}

	err = MergeChunks(tempDir, finalFilePath, initMessage.TotalChunks)
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error merging chunks")
		return
	}

	err = os.RemoveAll(tempDir)
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error removing temp directory")
		return
	}

	_, err = SaveFiles([]string{initMessage.FileName}, initMessage.UserID, db, usersRepository)
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error saving file to database")
		return
	}

	sockets.CloseWebSocket(ws, websocket.CloseNormalClosure, "file uploaded")
}

func HandleFileDownload(c *gin.Context, db *sql.DB) {
	fileIDParam := c.Param("fileID")

	ws, err := sockets.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error upgrading connection",
			"error":   err,
		})
		return
	}

	fileID, err := strconv.Atoi(fileIDParam)
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error converting fileID to int")
		return
	}

	fileRecord, err := GetFileByID(fileID, db)
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error getting file record")
		return
	}

	uploadsDir, err := path.CreatePathFromRoot("data/uploads")
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error creating uploads directory path")
		return
	}

	file, err := os.Open(fmt.Sprintf("%s/%s", uploadsDir, fileRecord.Name))
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error opening file")
		return
	}

	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error getting file stats")
		return
	}

	totalChunks := (stats.Size() + ChunkSize - 1) / ChunkSize

	err = ws.WriteJSON(map[string]any{
		"fileName":    filepath.Base(file.Name()),
		"totalChunks": totalChunks,
	})
	if err != nil {
		log.Println(err)
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error sending file info")
		return
	}

	buffer := make([]byte, ChunkSize)
	for {
		n, err := file.Read(buffer)

		if err != nil {
			if err == io.EOF {
				break
			}
		}

		err = ws.WriteMessage(websocket.BinaryMessage, buffer[:n])
		if err != nil {
			sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error sending file")
		}
	}

	sockets.CloseWebSocket(ws, websocket.CloseNormalClosure, "file sent")
}

func HandleDeleteFile(c *gin.Context, db *sql.DB) {
	fileID := c.Param("fileID")

	id, err := strconv.Atoi(fileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid file ID",
		})
		return
	}

	file, err := GetFileByID(id, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error getting file record",
			"error":   err,
		})
		return
	}

	finalFilePath, err := path.CreatePathFromRoot("data/uploads/" + file.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating final file path",
			"error":   err,
		})
		return
	}

	err = DeleteFile(id, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error deleting file",
			"error":   err,
		})
		return
	}

	err = os.Remove(finalFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error deleting file from disk",
			"error":   err,
		})
		return
	}

	c.Status(http.StatusNoContent)
}
