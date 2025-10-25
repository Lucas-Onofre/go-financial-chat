package errors

import "net/http"

var (
	ErrBadRequest    = NewAppError("BAD_REQUEST", "Invalid request data", http.StatusBadRequest, nil)
	ErrNotFound      = NewAppError("NOT_FOUND", "Resource not found", http.StatusNotFound, nil)
	ErrInternal      = NewAppError("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError, nil)
	ErrUnauthorized  = NewAppError("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized, nil)
	ErrUnprocessable = NewAppError("UNPROCESSABLE", "Unprocessable entity", http.StatusUnprocessableEntity, nil)
)

func NewAppError(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

func Wrap(base *AppError, err error) *AppError {
	if base == nil {
		return ErrInternal
	}
	return &AppError{
		Code:    base.Code,
		Message: base.Message,
		Status:  base.Status,
		Err:     err,
	}
}
