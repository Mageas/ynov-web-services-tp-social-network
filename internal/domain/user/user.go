package user

import "time"

// User represents a user in the system
type User struct {
	Email        string
	PasswordHash string
	Salt         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser creates a new User instance
func NewUser(email, passwordHash, salt string) *User {
	now := time.Now()
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		Salt:         salt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
