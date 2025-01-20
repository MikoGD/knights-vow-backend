package uploads

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"knights-vow/internal/database"
	"knights-vow/internal/resources/users"
	"knights-vow/pkg/path"
)

type File struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	CreatedDate   string `json:"createdDate"`
	OwnerID       string `json:"ownerID"`
	OwnerUsername string `json:"ownerUsername"`
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

func GetAllFiles() ([]File, error) {
	var filePath string
	var err error
	var rows *sql.Rows

	filePath, err = path.CreatePathFromRoot("internal/resources/uploads/sql/get-files-count.sql")

	if err != nil {
		return nil, err
	}

	rows = database.ExecuteSQLQuery(filePath)
	filesCount := 0

	if !rows.Next() {
		return nil, errors.New("no rows returned")
	}

	rows.Scan(&filesCount)

	filePath, err = path.CreatePathFromRoot("internal/resources/uploads/sql/get-all-files.sql")

	if err != nil {
		return nil, err
	}

	rows = database.ExecuteSQLQuery(filePath)

	files := make([]File, filesCount)

	i := 0
	for rows.Next() {
		file := File{}

		rows.Scan(&file.ID, &file.OwnerID, &file.Name, &file.CreatedDate, &file.OwnerUsername)
		files[i] = file
		i++
	}

	return files, nil
}
