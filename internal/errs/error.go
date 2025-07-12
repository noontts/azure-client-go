package errs

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"-"` // HTTP status code, not included in JSON
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

var (
	UnableToProceed = &AppError{Status: 5000, Message: "error unable to proceed", Code: http.StatusBadRequest}
	BadRequest      = &AppError{Status: 4000, Message: "bad request", Code: http.StatusBadRequest}
	NotFound        = &AppError{Status: 4040, Message: "not found", Code: http.StatusNotFound}
	// Add more custom errors here as needed
)
