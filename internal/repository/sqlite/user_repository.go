package sqlite

import (
	"context"
	"errors"
	"strings"

	"ynov-social-api/internal/domain/user"
	"ynov-social-api/internal/pkg/apperrors"

	"gorm.io/gorm"
)

// UserRepository implements user.Repository interface
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	model := &userModel{
		Email:        strings.ToLower(u.Email),
		PasswordHash: u.PasswordHash,
		Salt:         u.Salt,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return apperrors.ErrUserAlreadyExists
		}
		return apperrors.Wrap(err, 500, "failed to create user")
	}

	return nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var model userModel
	err := r.db.WithContext(ctx).
		Where("LOWER(email) = LOWER(?)", email).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.Wrap(err, 500, "failed to get user")
	}

	return &user.User{
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		Salt:         model.Salt,
	}, nil
}

// Exists checks if a user exists by email
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&userModel{}).
		Where("LOWER(email) = LOWER(?)", email).
		Count(&count).Error

	if err != nil {
		return false, apperrors.Wrap(err, 500, "failed to check user existence")
	}

	return count > 0, nil
}
