package users

import (
	"knights-vow/internal/database"
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

const (
	pathFromRoot = "internal/resources/users/sql"
)

// Save the user to the database.
func SaveUser(username string, password string) (int, error) {
	insertUserQuery, err := database.GetQuery(pathFromRoot + "/insert-user.sql")
	if err != nil {
		return -1, err
	}

	tx, err := database.Pool.Begin()
	if err != nil {
		return -1, err
	}

	result, err := tx.Exec(insertUserQuery, username, password)
	if err != nil {
		database.RollbackTx(tx)
		return -1, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		database.RollbackTx(tx)
		return -1, err
	}

	database.CommitTx(tx)

	return int(userID), nil
}

// Gets a user from the database by username. Returns nil, nil if the user does not exist.
func GetUserByUsername(username string) (*User, error) {
	selectUserByUsernameQuery, err := database.GetQuery(pathFromRoot + "/get-user-by-username.sql")
	if err != nil {
		return nil, err
	}

	rows, err := database.Pool.Query(selectUserByUsernameQuery, username)

	if err != nil {
		return nil, err
	}

	hasRow := rows.Next()

	if !hasRow && rows.Err() != nil {
		database.CloseRows(rows)
		return nil, rows.Err()
	}

	if !hasRow {
		database.CloseRows(rows)
		return nil, nil
	}

	user := &User{}
	err = rows.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		database.CloseRows(rows)
		return nil, err
	}

	return user, nil
}

// Gets a user from the database by username. Returns nil, nil if the user does not exist.
func GetUserByID(userID int) (*User, error) {
	selectUserByIDQuery, err := database.GetQuery(pathFromRoot + "/select-user-by-id.sql")
	if err != nil {
		return nil, err
	}

	rows, err := database.Pool.Query(selectUserByIDQuery, userID)
	if err != nil {
		database.CloseRows(rows)
		return nil, err
	}

	hasRow := rows.Next()

	if !hasRow && rows.Err() != nil {
		database.CloseRows(rows)
		return nil, rows.Err()
	}

	if !hasRow {
		database.CloseRows(rows)
		return nil, nil
	}

	user := &User{}
	err = rows.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		database.CloseRows(rows)
		return nil, err
	}

	return user, nil
}
