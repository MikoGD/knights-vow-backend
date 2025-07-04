package files

import (
	"encoding/json"
	"fmt"
	"io"
	"knights-vow/pkg/sockets"
	"math"
	"os"

	"github.com/gorilla/websocket"
	"github.com/mikogd/hextech/path"
)

type FilesService interface {
	GetFiles(fileName string) ([]File, error)
	UploadFile(ws *websocket.Conn)
	SaveFileFromUpload(fileName string, ownerID int, body io.ReadCloser) error
	GetFileForDownload(fileID int) (*os.File, error)
	DeleteFile(fileID int) error
}

type filesService struct {
	filesRepo FilesRepository
}

func createNewFilesService(filesRepo FilesRepository) FilesService {
	return &filesService{filesRepo}
}

func (s *filesService) GetFiles(fileName string) ([]File, error) {
	if fileName == "" {
		return s.filesRepo.GetAllFiles()
	}

	return s.filesRepo.GetFilesByName(fileName)
}

func (s *filesService) UploadFile(ws *websocket.Conn) {
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

		err = ws.WriteJSON(FileChunkSaveMessage{
			message:          "Chunk saved",
			chunkNumber:      i,
			uploadPercentage: int(percentageUploaded),
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

	_, err = s.filesRepo.SaveFiles([]string{initMessage.FileName}, initMessage.UserID)
	if err != nil {
		sockets.CloseWebSocket(ws, websocket.CloseInternalServerErr, "error saving file to database")
		return
	}

	sockets.CloseWebSocket(ws, websocket.CloseNormalClosure, "file uploaded")
}

func (s *filesService) SaveFileFromUpload(fileName string, ownerID int, body io.ReadCloser) error {
	outputPath, err := path.CreatePathFromRoot("data/uploads/" + fileName)
	if err != nil {
		return err
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer output.Close()
	defer body.Close()

	if _, err = io.Copy(output, body); err != nil {
		return err
	}

	_, err = s.filesRepo.SaveFiles([]string{fileName}, ownerID)
	if err != nil {
		return err
	}

	return nil
}

func (s *filesService) GetFileForDownload(fileID int) (*os.File, error) {
	fileRecord, err := s.filesRepo.GetFileByID(fileID)
	if err != nil {
		return nil, err
	}

	uploadsDir, err := path.CreatePathFromRoot("data/uploads")
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fmt.Sprintf("%s/%s", uploadsDir, fileRecord.Name))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *filesService) DeleteFile(fileID int) error {
	file, err := s.filesRepo.GetFileByID(fileID)
	if err != nil {
		return err
	}

	err = s.filesRepo.DeleteFile(fileID)
	if err != nil {
		return err
	}

	filePath, err := path.CreatePathFromRoot("data/uploads/" + file.Name)
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
