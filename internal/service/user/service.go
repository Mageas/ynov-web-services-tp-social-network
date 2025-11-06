package user

import (
	"context"
	"strings"

	"ynov-social-api/internal/domain/user"
	"ynov-social-api/internal/pkg/apperrors"
	"ynov-social-api/internal/pkg/validator"
	"ynov-social-api/internal/service/auth"
)

// Service handles user business logic
type Service struct {
	repo            user.Repository
	passwordService *auth.PasswordService
}

// NewService creates a new user service
func NewService(repo user.Repository, passwordService *auth.PasswordService) *Service {
	return &Service{
		repo:            repo,
		passwordService: passwordService,
	}
}

// Register registers a new user
func (s *Service) Register(ctx context.Context, email, password string) error {
	// Validate input
	v := validator.New()
	v.Required(email, "email")
	v.Email(email, "email")
	v.Required(password, "password")
	v.MinLength(password, 6, "password")

	if !v.Valid() {
		return apperrors.NewValidationError(v.GetErrors())
	}

	// Normalize email
	email = strings.ToLower(strings.TrimSpace(email))

	// Check if user already exists
	exists, err := s.repo.Exists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return apperrors.ErrUserAlreadyExists
	}

	// Hash password
	salt := s.passwordService.GenerateSalt(email)
	passwordHash := s.passwordService.HashPassword(password, salt)

	// Create user
	u := user.NewUser(email, passwordHash, salt)
	return s.repo.Create(ctx, u)
}

// Authenticate authenticates a user and returns their email
func (s *Service) Authenticate(ctx context.Context, email, password string) (string, error) {
	// Validate input
	v := validator.New()
	v.Required(email, "email")
	v.Required(password, "password")

	if !v.Valid() {
		return "", apperrors.NewValidationError(v.GetErrors())
	}

	// Normalize email
	email = strings.ToLower(strings.TrimSpace(email))

	// Get user
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	// Verify password
	if !s.passwordService.VerifyPassword(password, u.Salt, u.PasswordHash) {
		return "", apperrors.ErrInvalidCredentials
	}

	return u.Email, nil
}
