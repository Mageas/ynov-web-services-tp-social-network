package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password hashing and verification
type PasswordService struct {
	cost int // bcrypt cost factor
}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: bcrypt.DefaultCost,
	}
}

// HashPassword hashes a password using bcrypt
func (s *PasswordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a bcrypt hash
func (s *PasswordService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
