package user

import "time"

// User represents a user in the system
type User struct {
	Email        string
	PasswordHash string // bcrypt hash (salt is embedded in the hash)
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser creates a new User instance
func NewUser(email, passwordHash string) *User {
	now := time.Now()
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
