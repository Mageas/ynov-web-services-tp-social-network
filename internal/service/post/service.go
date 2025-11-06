package post

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"ynov-social-api/internal/domain/post"
	"ynov-social-api/internal/pkg/apperrors"
	"ynov-social-api/internal/pkg/validator"
)

// Service handles post business logic
type Service struct {
	repo post.Repository
}

// NewService creates a new post service
func NewService(repo post.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreatePost creates a new post
func (s *Service) CreatePost(ctx context.Context, author, content string) (*post.Post, error) {
	// Validate input
	content = strings.TrimSpace(content)
	v := validator.New()
	v.Required(content, "content")
	v.MaxLength(content, 400, "content")

	if !v.Valid() {
		return nil, apperrors.NewValidationError(v.GetErrors())
	}

	// Generate unique ID
	id, err := generateID()
	if err != nil {
		return nil, apperrors.Wrap(err, 500, "failed to generate post ID")
	}

	// Create post
	p := post.NewPost(id, author, content)
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	// Set likes count to 0 for new post
	p.LikesCount = 0
	return p, nil
}

// ListPosts retrieves posts with pagination
func (s *Service) ListPosts(ctx context.Context, beforeTimestamp int64, page, limit int) ([]*post.Post, error) {
	// If no beforeTimestamp provided, use current time + 1
	if beforeTimestamp <= 0 {
		beforeTimestamp = time.Now().Unix() + 1
	}

	return s.repo.ListBefore(ctx, beforeTimestamp, page, limit)
}

// LikePost adds a like to a post and returns the updated likes count
func (s *Service) LikePost(ctx context.Context, userEmail, postID string) (int, error) {
	if err := s.repo.AddLike(ctx, userEmail, postID); err != nil {
		return 0, err
	}

	// Retrieve the post to get the updated likes count
	p, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return 0, err
	}

	return p.LikesCount, nil
}

// UnlikePost removes a like from a post and returns the updated likes count
func (s *Service) UnlikePost(ctx context.Context, userEmail, postID string) (int, error) {
	if err := s.repo.RemoveLike(ctx, userEmail, postID); err != nil {
		return 0, err
	}

	// Retrieve the post to get the updated likes count
	p, err := s.repo.GetByID(ctx, postID)
	if err != nil {
		return 0, err
	}

	return p.LikesCount, nil
}

// generateID generates a unique ID for a post
func generateID() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
