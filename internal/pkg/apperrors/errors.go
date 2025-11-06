package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int
	Message string
	Err     error
	// ValidationErrors contains field-specific validation errors
	ValidationErrors map[string]string
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewValidationError creates a new validation error with field-specific errors
func NewValidationError(validationErrors map[string]string) *AppError {
	return &AppError{
		Code:             http.StatusBadRequest,
		Message:          "validation failed",
		ValidationErrors: validationErrors,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common application errors
var (
	ErrBadRequest         = New(http.StatusBadRequest, "bad request")
	ErrUnauthorized       = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden          = New(http.StatusForbidden, "forbidden")
	ErrNotFound           = New(http.StatusNotFound, "not found")
	ErrConflict           = New(http.StatusConflict, "conflict")
	ErrInternalServer     = New(http.StatusInternalServerError, "internal server error")
	ErrInvalidCredentials = New(http.StatusUnauthorized, "invalid credentials")
	ErrUserAlreadyExists  = New(http.StatusConflict, "user already exists")
	ErrPostNotFound       = New(http.StatusNotFound, "post not found")
	ErrInvalidToken       = New(http.StatusUnauthorized, "invalid token")
	ErrMissingAuth        = New(http.StatusUnauthorized, "missing authorization header")
)

// AsAppError converts an error to AppError if possible
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
