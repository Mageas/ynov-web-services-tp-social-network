package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"ynov-social-api/internal/types"

	"gorm.io/gorm"
)

// GORM-backed PostStore
type SQLitePostStore struct{ db *gorm.DB }

func (s *SQLiteStores) Posts() *SQLitePostStore { return &SQLitePostStore{db: s.DB} }

func (p *SQLitePostStore) Create(author, content string) (types.Post, error) {
	// generate id compatible with previous behavior
	buf := make([]byte, 12)
	_, _ = rand.Read(buf)
	id := hex.EncodeToString(buf)
	rec := postModel{ID: id, Author: author, Content: content, CreatedAt: time.Now().Unix()}
	if err := p.db.Create(&rec).Error; err != nil {
		return types.Post{}, err
	}
	return types.Post{ID: rec.ID, Author: rec.Author, Content: rec.Content, CreatedAt: rec.CreatedAt}, nil
}

func (p *SQLitePostStore) List(page, limit int) ([]types.Post, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit
	var rows []struct {
		ID         string
		Author     string
		Content    string
		CreatedAt  int64 `gorm:"column:created_at"`
		LikesCount int   `gorm:"column:likes_count"`
	}
	q := p.db.Table("posts").
		Select("posts.id, posts.author, posts.content, posts.created_at, COUNT(liked_posts.post_id) AS likes_count").
		Joins("LEFT JOIN liked_posts ON liked_posts.post_id = posts.id").
		Group("posts.id").
		Order("posts.created_at DESC, posts.id DESC").
		Offset(offset).Limit(limit)
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]types.Post, 0, len(rows))
	for _, r := range rows {
		out = append(out, types.Post{ID: r.ID, Author: r.Author, Content: r.Content, CreatedAt: r.CreatedAt, LikesCount: r.LikesCount})
	}
	return out, nil
}

// ListBefore implements cursor-based pagination using created_at as a strict cursor.
func (p *SQLitePostStore) ListBefore(beforeTs int64, page, limit int) ([]types.Post, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit
	var rows []struct {
		ID         string
		Author     string
		Content    string
		CreatedAt  int64 `gorm:"column:created_at"`
		LikesCount int   `gorm:"column:likes_count"`
	}
	q := p.db.Table("posts").
		Select("posts.id, posts.author, posts.content, posts.created_at, COUNT(liked_posts.post_id) AS likes_count").
		Joins("LEFT JOIN liked_posts ON liked_posts.post_id = posts.id").
		Group("posts.id").
		Order("posts.created_at DESC, posts.id DESC").
		Offset(offset).Limit(limit)

	if beforeTs > 0 {
		// Strictly older than the cursor: created_at < beforeTs
		q = q.Where("posts.created_at < ?", beforeTs)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]types.Post, 0, len(rows))
	for _, r := range rows {
		out = append(out, types.Post{ID: r.ID, Author: r.Author, Content: r.Content, CreatedAt: r.CreatedAt, LikesCount: r.LikesCount})
	}
	return out, nil
}

func (p *SQLitePostStore) Like(userEmail, postID string) error {
	// ensure post exists
	var post postModel
	if err := p.db.First(&post, "id = ?", postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPostNotFound
		}
		return err
	}
	rec := likedPostModel{UserEmail: userEmail, PostID: postID}
	return p.db.FirstOrCreate(&rec, likedPostModel{UserEmail: userEmail, PostID: postID}).Error
}

func (p *SQLitePostStore) Unlike(userEmail, postID string) error {
	// ensure post exists
	var post postModel
	if err := p.db.First(&post, "id = ?", postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPostNotFound
		}
		return err
	}
	return p.db.Where("user_email = ? AND post_id = ?", userEmail, postID).Delete(&likedPostModel{}).Error
}
