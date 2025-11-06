package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

// PasswordService handles password hashing and verification
type PasswordService struct{}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// HashPassword hashes a password with the given salt
func (s *PasswordService) HashPassword(password, salt string) string {
	saltBytes, _ := hex.DecodeString(salt)
	sum := sha256.Sum256(append([]byte(password), saltBytes...))
	return hex.EncodeToString(sum[:])
}

// GenerateSalt generates a salt from an email (deterministic)
func (s *PasswordService) GenerateSalt(email string) string {
	return hex.EncodeToString([]byte(email))
}

// VerifyPassword verifies a password against a hash
func (s *PasswordService) VerifyPassword(password, salt, hash string) bool {
	return s.HashPassword(password, salt) == hash
}
