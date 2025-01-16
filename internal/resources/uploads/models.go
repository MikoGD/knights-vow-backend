package uploads

import (
	"errors"
	"log"
	"time"

	"knights-vow/internal/database"
	"knights-vow/internal/resources/users"
	"knights-vow/pkg/path"
)

type File struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CreatedDate string `json:"created_date"`
	Owner       string `json:"owner"`
}

func createArgs(args ...any) []any {
	return args
}

func SaveFiles(fileNames []string, ownerName string) (int, error) {
	if len(fileNames) == 0 {
		return 0, errors.New("no files to save")
	}

	user, err := users.GetUserByUsername(ownerName)

	if err != nil {
		log.Fatalf("Error getting user by username: %v", err)
	}

	filePath, err := path.CreatePathFromRoot("internal/resources/uploads/sql/add-file.sql")

	if err != nil {
		return 0, err
	}

	placeholders := make([][]any, 0, len(fileNames))

	for _, fileName := range fileNames {
		placeholders = append(placeholders, createArgs(fileName, user.ID, time.Now().Format(time.RFC3339)))
	}

	results := database.ExecuteSQLStatementWithMultipleArgs(filePath, placeholders)

	totalRowsInserted := len(results)

	return totalRowsInserted, nil
}
