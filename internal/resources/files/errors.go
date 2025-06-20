package files

import (
	"errors"
	"net/http"
)

type NoFilesError struct{}

func (e *NoFilesError) Error() string {
	return "Error no files found to save"
}

type ReadFilesDatabaseError struct{}

func (e *ReadFilesDatabaseError) Error() string {
	return "Error failed to read files database"
}

type errorResponse struct {
	message string
	err     *error
}

func createErrorResponse(err error) (int, *errorResponse) {
	if errors.Is(err, &NoFilesError{}) {
		return http.StatusNotFound, &errorResponse{
			message: err.Error(),
		}
	}

	if errors.Is(err, &ReadFilesDatabaseError{}) {
		return http.StatusConflict, &errorResponse{
			message: err.Error(),
		}
	}

	return http.StatusInternalServerError, &errorResponse{
		message: "something went wrong",
		err:     &err,
	}
}

func createErrorResponseWithCustomDefaultMesage(err error, message string) (int, *errorResponse) {
	if errors.Is(err, &NoFilesError{}) {
		return http.StatusNotFound, &errorResponse{
			message: err.Error(),
		}
	}

	if errors.Is(err, &ReadFilesDatabaseError{}) {
		return http.StatusConflict, &errorResponse{
			message: err.Error(),
		}
	}

	return http.StatusInternalServerError, &errorResponse{
		message: message,
		err:     &err,
	}
}
