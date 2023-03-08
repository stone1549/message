package service

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (er ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(er.Status)

	return nil
}

func (er ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", er.Status, er.Message)
}

func NewNotFoundErr(message string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func NewInternalServerErr(message string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

func NewBadRequestErr(message string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func NewUnauthorizedErr(message string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}
