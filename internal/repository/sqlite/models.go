package sqlite

// userModel represents the database model for users
type userModel struct {
	Email        string `gorm:"primaryKey"`
	PasswordHash string
	Salt         string
}

// TableName overrides the table name
func (userModel) TableName() string {
	return "users"
}

// postModel represents the database model for posts
type postModel struct {
	ID        string `gorm:"primaryKey"`
	Author    string `gorm:"index"`
	Content   string
	CreatedAt int64 `gorm:"index"`
}

// TableName overrides the table name
func (postModel) TableName() string {
	return "posts"
}

// likeModel represents the database model for likes
type likeModel struct {
	UserEmail string `gorm:"primaryKey;column:user_email"`
	PostID    string `gorm:"primaryKey;column:post_id;index"`
}

// TableName overrides the table name
func (likeModel) TableName() string {
	return "liked_posts"
}
