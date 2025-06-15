package users

import (
	"database/sql"
	"knights-vow/internal/database"
	"log"
	"os"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var (
	repository UserRepository
	db         *sql.DB
)

func TestMain(m *testing.M) {
	originalWD, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get cwd: %v\n", err)
	}

	pathToSystemHome := os.Getenv("HOME")
	if pathToSystemHome == "" {
		log.Fatalln("Failed to get HOME environment variable")
	}

	projectRootPath := path.Join(pathToSystemHome, "workspace/projects/knights-vow/backend")
	err = os.Chdir(projectRootPath)
	if err != nil {
		log.Fatalf("Failed to change directory: %v\n", err)
	}

	defer os.Chdir(originalWD)

	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to create in-memory database %s\n", err)
	}
	defer db.Close()

	database.CreateTables(db)

	repository = CreateNewUserRepository(db)

	os.Exit(m.Run())
}

func TestSaveUser(t *testing.T) {
	username := "test-username-1"
	password := "test-password-123"

	userID, err := repository.SaveUser(username, password)
	require.NoError(t, err)

	expectedUser := &User{ID: 1, Username: username, Password: password}

	user, err := repository.GetUserByID(userID)
	require.Equal(t, expectedUser, user)
}

func TestGetUserByUsername(t *testing.T) {
	_, err := db.Exec("INSERT INTO Users (username, password) VALUES ('test-username-2', 'test-password-2')")
	if err != nil {
		log.Fatalf("Failed to add test row to database: %v\n", err)
	}

	expectedUser := &User{ID: 2, Username: "test-username-2", Password: "test-password-2"}
	user, err := repository.GetUserByUsername("test-username-2")
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
}

func TestGetUserByUserID(t *testing.T) {
	_, err := db.Exec("INSERT INTO Users (username, password) VALUES ('test-username-3', 'test-password-3')")
	if err != nil {
		log.Fatalf("Failed to add test row to database: %v\n", err)
	}

	expectedUser := &User{ID: 3, Username: "test-username-3", Password: "test-password-3"}
	user, err := repository.GetUserByID(3)
	require.NoError(t, err)
	require.Equal(t, expectedUser, user)
}
