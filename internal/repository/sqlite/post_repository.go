package sqlite

import (
	"context"
	"errors"
	"time"

	"ynov-social-api/internal/domain/post"
	"ynov-social-api/internal/pkg/apperrors"

	"gorm.io/gorm"
)

// PostRepository implements post.Repository interface
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new PostRepository
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create creates a new post
func (r *PostRepository) Create(ctx context.Context, p *post.Post) error {
	model := &postModel{
		ID:        p.ID,
		UserEmail: p.Author, // Author is the user email
		Content:   p.Content,
		CreatedAt: p.CreatedAt.Unix(),
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return apperrors.Wrap(err, 500, "failed to create post")
	}

	return nil
}

// GetByID retrieves a post by ID
func (r *PostRepository) GetByID(ctx context.Context, id string) (*post.Post, error) {
	var model postModel
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrPostNotFound
		}
		return nil, apperrors.Wrap(err, 500, "failed to get post")
	}

	// Get likes count
	var likesCount int64
	r.db.WithContext(ctx).
		Model(&likeModel{}).
		Where("post_id = ?", id).
		Count(&likesCount)

	return &post.Post{
		ID:         model.ID,
		Author:     model.UserEmail, // UserEmail is the author email
		Content:    model.Content,
		CreatedAt:  time.Unix(model.CreatedAt, 0),
		LikesCount: int(likesCount),
	}, nil
}

// ListBefore retrieves posts created before a given timestamp with pagination
func (r *PostRepository) ListBefore(ctx context.Context, beforeTimestamp int64, page, limit int) ([]*post.Post, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	var results []struct {
		ID         string
		UserEmail  string
		Content    string
		CreatedAt  int64
		LikesCount int64
	}

	query := r.db.WithContext(ctx).
		Table("posts").
		Select("posts.id, posts.user_email, posts.content, posts.created_at, COUNT(liked_posts.post_id) AS likes_count").
		Joins("LEFT JOIN liked_posts ON liked_posts.post_id = posts.id").
		Group("posts.id").
		Order("posts.created_at DESC, posts.id DESC")

	if beforeTimestamp > 0 {
		query = query.Where("posts.created_at < ?", beforeTimestamp)
	}

	err := query.Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return nil, apperrors.Wrap(err, 500, "failed to list posts")
	}

	posts := make([]*post.Post, 0, len(results))
	for _, r := range results {
		posts = append(posts, &post.Post{
			ID:         r.ID,
			Author:     r.UserEmail, // UserEmail is the author email
			Content:    r.Content,
			CreatedAt:  time.Unix(r.CreatedAt, 0),
			LikesCount: int(r.LikesCount),
		})
	}

	return posts, nil
}

// AddLike adds a like to a post
func (r *PostRepository) AddLike(ctx context.Context, userEmail, postID string) error {
	// Check if post exists
	exists, err := r.Exists(ctx, postID)
	if err != nil {
		return err
	}
	if !exists {
		return apperrors.ErrPostNotFound
	}

	like := &likeModel{
		UserEmail: userEmail,
		PostID:    postID,
	}

	// Use FirstOrCreate to avoid duplicate likes
	err = r.db.WithContext(ctx).
		Where(likeModel{UserEmail: userEmail, PostID: postID}).
		FirstOrCreate(like).Error

	if err != nil {
		return apperrors.Wrap(err, 500, "failed to add like")
	}

	return nil
}

// RemoveLike removes a like from a post
func (r *PostRepository) RemoveLike(ctx context.Context, userEmail, postID string) error {
	// Check if post exists
	exists, err := r.Exists(ctx, postID)
	if err != nil {
		return err
	}
	if !exists {
		return apperrors.ErrPostNotFound
	}

	err = r.db.WithContext(ctx).
		Where("user_email = ? AND post_id = ?", userEmail, postID).
		Delete(&likeModel{}).Error

	if err != nil {
		return apperrors.Wrap(err, 500, "failed to remove like")
	}

	return nil
}

// Exists checks if a post exists by ID
func (r *PostRepository) Exists(ctx context.Context, postID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&postModel{}).
		Where("id = ?", postID).
		Count(&count).Error

	if err != nil {
		return false, apperrors.Wrap(err, 500, "failed to check post existence")
	}

	return count > 0, nil
}
