package files

import (
	"fmt"
	"io"
	"knights-vow/pkg/path"
	"os"
	"path/filepath"
	"strconv"
)

type FileUploadInitMessage struct {
	FileName    string `json:"fileName"`
	TotalChunks int    `json:"totalChunks"`
	UserID      int    `json:"userID"`
}

func CreateTempDir(userID int, fileName string) (string, error) {
	tempDir, err := path.CreatePathFromRoot(fmt.Sprintf("%s/%d/%s", PathToTempDirectory, userID, fileName))

	if err != nil {
		return "", err
	}

	err = os.MkdirAll(tempDir, os.ModePerm)

	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func SaveChunk(tempDir string, chunkNumber int, data []byte) error {
	chunkPath := filepath.Join(tempDir, strconv.Itoa(chunkNumber))

	chunk, err := os.Create(chunkPath)

	if err != nil {
		return err
	}
	_, err = chunk.Write(data)

	if err != nil {
		return err
	}

	return nil
}

func MergeChunks(chunksTempDir string, fileDestPath string, totalChunks int) error {
	finalFile, err := os.OpenFile(fileDestPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)

	if err != nil {
		return err
	}

	defer finalFile.Close()

	for i := 1; i <= totalChunks; i++ {
		chunkPath := filepath.Join(chunksTempDir, strconv.Itoa(i))
		chunk, err := os.Open(chunkPath)

		if err != nil {
			return err
		}

		defer chunk.Close()

		_, err = io.Copy(finalFile, chunk)

		if err != nil {
			return err
		}
	}

	return nil
}
