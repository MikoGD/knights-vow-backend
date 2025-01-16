package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"knights-vow/pkg/path"
)

var Pool *sql.DB

// Runs the same SQL statement for multiple arguments. For example, multiple inserts.
func ExecuteSQLStatementWithMultipleArgs(statementFilePath string, args [][]any) []sql.Result {
	results := make([]sql.Result, len(args))

	tx, err := Pool.Begin()

	if err != nil {
		tx.Rollback()
		log.Fatalf("Error beginning transaction: %v", err)
	}

	content, err := os.ReadFile(statementFilePath)
	statement := string(content)
	statement = strings.TrimSpace(statement)

	for i, arg := range args {
		if err != nil {
			log.Fatalf("Error reading SQL file: %v", err)
		}

		result, err := tx.Exec(statement, arg...)

		if err != nil {
			tx.Rollback()
			log.Fatalf("Error executing \"%v\" statement: %v", statementFilePath, err)
		}

		results[i] = result
	}

	tx.Commit()

	return results
}

func ExecuteSQLStatement(statementFilePath string, args ...any) sql.Result {
	content, err := os.ReadFile(statementFilePath)

	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	statement := string(content)

	statement = strings.TrimSpace(statement)

	tx, err := Pool.Begin()

	if err != nil {
		tx.Rollback()
		log.Fatalf("Error beginning transaction: %v", err)
	}

	result, err := tx.Exec(statement, args...)

	if err != nil {
		tx.Rollback()
		log.Fatalf("Error executing \"%v\" statement: %v", statementFilePath, err)
	}

	tx.Commit()

	return result
}

func ExecuteSQLQuery(queryFilePath string, args ...any) *sql.Rows {
	content, err := os.ReadFile(queryFilePath)

	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	query := strings.TrimSpace(string(content))

	stmt, err := Pool.Prepare(query)

	if err != nil {
		log.Fatalf("Error preparing \"%v\" query: %v", queryFilePath, err)
	}

	defer stmt.Close()

	rows, err := stmt.Query(args...)

	if err != nil {
		log.Fatalf("Error executing \"%v\" query: %v", queryFilePath, err)
	}

	return rows
}

func createTables() {
	createUsersTableSQL, err := path.CreatePathFromRoot("internal/database/sql/create-users-table.sql")

	if err != nil {
		log.Fatalf("Error creating path from root: %v", err)
	}

	createFilesTableSQL, err := path.CreatePathFromRoot("internal/database/sql/create-files-table.sql")

	if err != nil {
		log.Fatalf("Error creating path from root: %v", err)
	}

	ExecuteSQLStatement(createUsersTableSQL)
	ExecuteSQLStatement(createFilesTableSQL)
}

func InitDatabase() {
	var err error

	databasePath, err := path.CreatePathFromRoot("/data/databases/knights-vow.db")

	if err != nil {
		log.Fatalf("Error creating path from root: %v", err)
	}

	Pool, err = sql.Open("sqlite3", databasePath)

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	createTables()
}

func CloseDatabase() {
	Pool.Close()
}
