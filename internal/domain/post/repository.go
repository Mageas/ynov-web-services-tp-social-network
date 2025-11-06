package post

import "context"

// Repository defines the interface for post data access
type Repository interface {
	// Create creates a new post
	Create(ctx context.Context, post *Post) error

	// GetByID retrieves a post by ID
	GetByID(ctx context.Context, id string) (*Post, error)

	// ListBefore retrieves posts created before a given timestamp with pagination
	ListBefore(ctx context.Context, beforeTimestamp int64, page, limit int) ([]*Post, error)

	// AddLike adds a like to a post
	AddLike(ctx context.Context, userEmail, postID string) error

	// RemoveLike removes a like from a post
	RemoveLike(ctx context.Context, userEmail, postID string) error

	// Exists checks if a post with the given ID exists
	Exists(ctx context.Context, postID string) (bool, error)
}
