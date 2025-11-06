package response

import (
	"encoding/json"
	"net/http"

	"ynov-social-api/internal/api/dto"
	"ynov-social-api/internal/pkg/apperrors"
)

// JSON writes a JSON success response
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Error writes a JSON error response
func Error(w http.ResponseWriter, err error) {
	appErr, ok := apperrors.AsAppError(err)
	if !ok {
		appErr = apperrors.ErrInternalServer
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)

	response := dto.ErrorResponse{
		Status:           appErr.Code,
		Message:          appErr.Message,
		ValidationErrors: appErr.ValidationErrors,
	}

	json.NewEncoder(w).Encode(response)
}

// Created writes a 201 Created response (for POST requests)
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

// OK writes a 200 OK response
func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

// NoContent writes a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
