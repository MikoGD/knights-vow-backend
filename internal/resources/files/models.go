package files

import (
	"errors"
	"time"

	"knights-vow/internal/database"
	"knights-vow/internal/resources/users"
)

type File struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	CreatedDate   string `json:"createdDate"`
	OwnerID       int    `json:"ownerID"`
	OwnerUsername string `json:"ownerUsername"`
}

const (
	pathFromRoot = "internal/resources/files/sql"
)

func SaveFiles(fileNames []string, ownerID int) (int, error) {
	if len(fileNames) == 0 {
		return 0, errors.New("no files to save")
	}

	user, err := users.GetUserByID(ownerID)
	if err != nil {
		return 0, err
	}

	tx, err := database.Pool.Begin()
	if err != nil {
		return 0, err
	}

	addFileQuery, err := database.GetQuery(pathFromRoot + "/insert-file.sql")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	filesSaved := 0
	for _, fileName := range fileNames {
		_, err = tx.Exec(addFileQuery, fileName, user.ID, time.Now().Format(time.RFC3339))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		filesSaved++
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return filesSaved, nil
}

func GetAllFilesCount() (int, error) {
	getFileCountQuery, err := database.GetQuery(pathFromRoot + "/select-files-count.sql")
	if err != nil {
		return 0, err
	}

	rows, err := database.Pool.Query(getFileCountQuery)
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, errors.New("no rows returned")
	}

	filesCount := 0
	rows.Scan(&filesCount)

	database.CloseRows(rows)

	return filesCount, nil
}

func GetAllFiles() ([]File, error) {
	filesCount, err := GetAllFilesCount()
	if err != nil {
		return nil, err
	}

	getAllFilesQuery, err := database.GetQuery(pathFromRoot + "/select-all-files.sql")
	if err != nil {
		return nil, err
	}

	stmt, err := database.Pool.Prepare(getAllFilesQuery)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		database.CloseStmt(stmt)
		return nil, err
	}

	files := make([]File, filesCount)
	i := 0
	for rows.Next() {
		file := File{}

		rows.Scan(&file.ID, &file.OwnerID, &file.Name, &file.CreatedDate, &file.OwnerUsername)
		files[i] = file
		i++
	}

	database.CloseRows(rows)
	database.CloseStmt(stmt)

	return files, nil
}

func GetFileByID(fileID int) (*File, error) {
	file := &File{}

	selectFileByIDQuery, err := database.GetQuery(pathFromRoot + "/select-file-by-id.sql")
	if err != nil {
		return nil, err
	}

	rows, err := database.Pool.Query(selectFileByIDQuery, fileID)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	err = rows.Scan(&file.ID, &file.Name, &file.CreatedDate, &file.OwnerID, &file.OwnerUsername)
	if err != nil {
		database.CloseRows(rows)
		return nil, err
	}

	database.CloseRows(rows)

	return file, nil
}
