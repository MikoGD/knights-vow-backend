package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"knights-vow/pkg/path"
)

const (
	pathFromRoot = "internal/database/sql"
)

func CreateTables(db *sql.DB) {
	createUsersTableQuery, err := GetQuery(pathFromRoot + "/create-users-table.sql")
	if err != nil {
		log.Fatalf("Error getting query: %v", err)
	}

	createFilesTableQuery, err := GetQuery(pathFromRoot + "/create-files-table.sql")
	if err != nil {
		log.Fatalf("Error getting query: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Error beginning transaction: %v", err)
	}

	_, err = tx.Exec(createUsersTableQuery)
	if err != nil {
		RollbackTx(tx)
		log.Fatalf("Error creating users table: %v", err)
	}

	_, err = tx.Exec(createFilesTableQuery)
	if err != nil {
		RollbackTx(tx)
		log.Fatalf("Error creating files table: %v", err)
	}

	CommitTx(tx)
}

func InitDatabase() *sql.DB {
	var err error

	databasePath, err := path.CreatePathFromRoot("/data/databases/knights-vow.db")

	if err != nil {
		log.Fatalf("Error creating path from root: %v", err)
	}

	db, err := sql.Open("sqlite3", databasePath)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	CreateTables(db)

	return db
}

func CloseDatabase(db *sql.DB) {
	db.Close()
}

// pathFromRoot is the relative path from the root, root is pre-prended in the function
func GetQuery(pathFromRoot string) (string, error) {
	queryFilePath, err := path.CreatePathFromRoot(pathFromRoot)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(queryFilePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func CloseRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Printf("Error closing rows: %v", err)
	}
}

func CloseStmt(stmt *sql.Stmt) {
	err := stmt.Close()
	if err != nil {
		log.Printf("Error closing statement: %v", err)
	}
}

func CommitTx(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
	}
}

func RollbackTx(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		log.Printf("Error rolling back transaction: %v", err)
	}
}
