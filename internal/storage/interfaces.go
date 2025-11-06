package storage

import (
	"errors"
	"ynov-social-api/internal/types"
)

// Domain errors
var ErrPostNotFound = errors.New("post not found")

type UserStore interface {
	CreateUser(email, password string) error
	VerifyCredentials(email, password string) bool
}

type PostStore interface {
	Create(author, content string) (types.Post, error)
	// List returns posts ordered by newest first with pagination.
	// page is 1-based; limit is the page size.
	List(page, limit int) ([]types.Post, error)
	// ListBefore returns posts strictly older than the beforeTs timestamp,
	// ordered by newest first. Supports page (1-based) and limit.
	ListBefore(beforeTs int64, page, limit int) ([]types.Post, error)
	Like(userEmail, postID string) error
	Unlike(userEmail, postID string) error
}
