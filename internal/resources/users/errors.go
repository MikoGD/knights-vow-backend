package users

import (
	"errors"
	"fmt"
	"net/http"
)

// For creating user
type UserExistsError struct {
	username string
}

func (e *UserExistsError) Error() string {
	return fmt.Sprintf("User with username %s", e.username)
}

// For logging in
type UserDoesNotExistError struct {
	username string
}

func (e *UserDoesNotExistError) Error() string {
	return fmt.Sprintf("Could not find user with username %s", e.username)
}

type InvalidLoginError struct{}

func (e *InvalidLoginError) Error() string {
	return "Invalid login username or password"
}

type InvalidTokenError struct{}

func (e *InvalidTokenError) Error() string {
	return "Invalid JWT"
}

type errorResponse struct {
	message string
	err     *error
}

func createErrorResponse(err error) (int, *errorResponse) {
	if errors.Is(err, &UserDoesNotExistError{}) {
		return http.StatusNotFound, &errorResponse{
			message: err.Error(),
		}
	}

	if errors.Is(err, &UserExistsError{}) {
		return http.StatusConflict, &errorResponse{
			message: err.Error(),
		}
	}

	if errors.Is(err, &InvalidLoginError{}) || errors.Is(err, &InvalidTokenError{}) {
		return http.StatusUnauthorized, &errorResponse{
			message: err.Error(),
		}
	}

	return http.StatusInternalServerError, &errorResponse{
		message: "something went wrong",
		err:     &err,
	}
}
