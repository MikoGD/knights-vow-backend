package path

import (
	"log"
	"os"
	"path/filepath"
)

// Creates a path from root contains the functions to create the path from the root of the project.
func CreatePathFromRoot(filePath string) (string, error) {
	rootPath, err := os.Getwd()

	if err != nil {
		log.Printf("Error getting current directory: %v", err)
		return "", err
	}

	pathFromRoot := filepath.Join(rootPath, filePath)

	return pathFromRoot, nil
}
