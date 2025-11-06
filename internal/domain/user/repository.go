package user

import "context"

// Repository defines the interface for user data access
type Repository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)

	// Exists checks if a user with the given email exists
	Exists(ctx context.Context, email string) (bool, error)
}
