package users

import (
	"knights-vow/internal/database"
	"knights-vow/pkg/path"
)

type UserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Save the user to the database.
func SaveUser(username string, password string) (int, error) {
	filePath, err := path.CreatePathFromRoot("internal/resources/users/sql/create-user.sql")

	if err != nil {
		return -1, err
	}

	result := database.ExecuteSQLStatement(filePath, username, password)

	userID, err := result.LastInsertId()

	if err != nil {
		return -1, err
	}

	return int(userID), nil
}

// Gets a user from the database by username. Returns nil, nil if the user does not exist.
func GetUserByUsername(username string) (*User, error) {
	filePath, err := path.CreatePathFromRoot("internal/resources/users/sql/get-user-by-username.sql")

	if err != nil {
		return nil, err
	}

	rows := database.ExecuteSQLQuery(filePath, username)

	defer rows.Close()

	hasRow := rows.Next()

	if !hasRow && rows.Err() != nil {
		return nil, rows.Err()
	}

	if !hasRow {
		return nil, nil
	}

	user := &User{}

	err = rows.Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Gets a user from the database by username. Returns nil, nil if the user does not exist.
func GetUserByID(userID int) (*User, error) {
	filePath, err := path.CreatePathFromRoot("internal/resources/users/sql/get-user-by-id.sql")

	if err != nil {
		return nil, err
	}

	rows := database.ExecuteSQLQuery(filePath, userID)

	defer rows.Close()

	hasRow := rows.Next()

	if !hasRow && rows.Err() != nil {
		return nil, rows.Err()
	}

	if !hasRow {
		return nil, nil
	}

	user := &User{}

	err = rows.Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return nil, err
	}

	return user, nil
}
