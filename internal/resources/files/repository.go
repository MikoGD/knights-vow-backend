package files

import (
	"database/sql"
	"fmt"
	"knights-vow/internal/database"
	"knights-vow/internal/resources/users"
	"time"
)

const (
	pathFromRoot = "internal/resources/files/sql"
)

type FilesRepository interface {
	SaveFiles(fileNames []string, ownerID int) (int, error)
	GetAllFilesCount() (int, error)
	GetAllFiles() ([]File, error)
	GetFileByID(fileID int) (*File, error)
	GetFilesByName(fileName string) ([]File, error)
	DeleteFile(fileID int) error
}

type filesRepository struct {
	db              *sql.DB
	usersRepository users.UserRepository
}

func CreateNewFilesRepository(db *sql.DB, usersRepository users.UserRepository) FilesRepository {
	return &filesRepository{db, usersRepository}
}

func (r *filesRepository) SaveFiles(fileNames []string, ownerID int) (int, error) {
	if len(fileNames) == 0 {
		return 0, &NoFilesError{}
	}

	user, err := r.usersRepository.GetUserByID(ownerID)
	if err != nil {
		return 0, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	addFileQuery, err := database.GetQuery(pathFromRoot + "/insert-file.sql")
	if err != nil {
		database.RollbackTx(tx)
		return 0, err
	}

	filesSaved := 0
	for _, fileName := range fileNames {
		_, err = tx.Exec(addFileQuery, fileName, user.ID, time.Now().Format(time.RFC3339))
		if err != nil {
			database.RollbackTx(tx)
			return 0, err
		}

		filesSaved++
	}

	database.CommitTx(tx)

	return filesSaved, nil
}

func (r *filesRepository) GetAllFilesCount() (int, error) {
	getFileCountQuery, err := database.GetQuery(pathFromRoot + "/select-files-count.sql")
	if err != nil {
		return 0, err
	}

	rows, err := r.db.Query(getFileCountQuery)
	if err != nil {
		return -1, err
	}

	if !rows.Next() {
		database.CloseRows(rows)
		return -1, &ReadFilesDatabaseError{}
	}

	filesCount := 0
	rows.Scan(&filesCount)

	database.CloseRows(rows)

	return filesCount, nil
}

func (r *filesRepository) GetAllFiles() ([]File, error) {
	filesCount, err := r.GetAllFilesCount()
	if err != nil {
		return nil, err
	}

	getAllFilesQuery, err := database.GetQuery(pathFromRoot + "/select-all-files.sql")
	if err != nil {
		return nil, err
	}

	stmt, err := r.db.Prepare(getAllFilesQuery)
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

func (r *filesRepository) GetFileByID(fileID int) (*File, error) {
	file := &File{}

	selectFileByIDQuery, err := database.GetQuery(pathFromRoot + "/select-file-by-id.sql")
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(selectFileByIDQuery, fileID)
	if err != nil {
		database.CloseRows(rows)
		return nil, err
	}

	if !rows.Next() {
		database.CloseRows(rows)
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

func (r *filesRepository) GetFilesByName(fileName string) ([]File, error) {
	files := make([]File, 0)

	selectFileByName, err := database.GetQuery(pathFromRoot + "/select-files-by-name.sql")
	if err != nil {
		return files, err
	}

	rows, err := r.db.Query(
		selectFileByName,
		fmt.Sprintf("%s%%", fileName),
		fmt.Sprintf("%%%s%%", fileName),
	)

	if err != nil {
		database.CloseRows(rows)
		return files, err
	}

	for rows.Next() {
		file := File{}

		if err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.CreatedDate,
			&file.OwnerID,
			&file.OwnerUsername,
		); err != nil {
			database.CloseRows(rows)
			return files, err
		}

		files = append(files, file)
	}

	if rows.Err() != nil {
		database.CloseRows(rows)
		return files, err
	}

	database.CloseRows(rows)

	return files, nil
}

func (r *filesRepository) DeleteFile(fileID int) error {
	deleteFileQuery, err := database.GetQuery(pathFromRoot + "/delete-file-by-id.sql")
	if err != nil {
		return err
	}

	_, err = r.db.Exec(deleteFileQuery, fileID)
	if err != nil {
		return err
	}

	return nil
}
