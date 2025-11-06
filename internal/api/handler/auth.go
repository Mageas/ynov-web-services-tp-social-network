package handler

import (
	"encoding/json"
	"net/http"

	"ynov-social-api/internal/api/dto"
	"ynov-social-api/internal/api/response"
	"ynov-social-api/internal/pkg/apperrors"
	"ynov-social-api/internal/pkg/logger"
	"ynov-social-api/internal/service/auth"
	"ynov-social-api/internal/service/user"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userService *user.Service
	jwtService  *auth.JWTService
	logger      *logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService *user.Service, jwtService *auth.JWTService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
		logger:      logger,
	}
}

// Signup handles user registration
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, apperrors.New(http.StatusMethodNotAllowed, "method not allowed"))
		return
	}

	var req dto.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperrors.New(http.StatusBadRequest, "invalid JSON"))
		return
	}

	if err := h.userService.Register(r.Context(), req.Email, req.Password); err != nil {
		h.logger.Error("Failed to register user: %v", err)
		response.Error(w, err)
		return
	}

	response.NoContent(w)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, apperrors.New(http.StatusMethodNotAllowed, "method not allowed"))
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperrors.New(http.StatusBadRequest, "invalid JSON"))
		return
	}

	email, err := h.userService.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to authenticate user: %v", err)
		response.Error(w, err)
		return
	}

	token, err := h.jwtService.GenerateToken(email)
	if err != nil {
		h.logger.Error("Failed to generate token: %v", err)
		response.Error(w, apperrors.ErrInternalServer)
		return
	}

	response.OK(w, dto.TokenResponse{Token: token})
}
