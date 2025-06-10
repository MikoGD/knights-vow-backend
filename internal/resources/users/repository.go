package users

import (
	"database/sql"
	"knights-vow/internal/database"
)

const (
	pathFromRoot = "internal/resources/users/sql"
)

type UserRepository interface {
	SaveUser(username string, password string) (int, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByID(userID int) (*User, error)
}

type userRepository struct {
	db *sql.DB
}

func CreateNewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

// Save the user to the database.
// Returns the user ID of user saved
func (r *userRepository) SaveUser(username string, password string) (int, error) {
	insertUserQuery, err := database.GetQuery(pathFromRoot + "/insert-user.sql")
	if err != nil {
		return -1, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		database.RollbackTx(tx)
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
func (r *userRepository) GetUserByUsername(username string) (*User, error) {
	selectUserByUsernameQuery, err := database.GetQuery(pathFromRoot + "/select-user-by-username.sql")
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(selectUserByUsernameQuery, username)

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

	database.CloseRows(rows)

	return user, nil
}

// Gets a user from the database by username. Returns nil, nil if the user does not exist.
func (r *userRepository) GetUserByID(userID int) (*User, error) {
	selectUserByIDQuery, err := database.GetQuery(pathFromRoot + "/select-user-by-id.sql")
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(selectUserByIDQuery, userID)
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

	database.CloseRows(rows)

	return user, nil
}
