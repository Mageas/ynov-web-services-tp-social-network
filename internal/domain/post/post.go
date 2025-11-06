package post

import "time"

// Post represents a post in the system
type Post struct {
	ID         string
	Author     string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LikesCount int
}

// NewPost creates a new Post instance
func NewPost(id, author, content string) *Post {
	now := time.Now()
	return &Post{
		ID:        id,
		Author:    author,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
