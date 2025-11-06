package sqlite

// userModel represents the database model for users
type userModel struct {
	Email        string `gorm:"primaryKey"`
	PasswordHash string // bcrypt hash (salt is embedded in the hash)
}

// TableName overrides the table name
func (userModel) TableName() string {
	return "users"
}

// postModel represents the database model for posts
type postModel struct {
	ID        string `gorm:"primaryKey"`
	UserEmail string `gorm:"column:user_email;index;not null"`
	Content   string
	CreatedAt int64 `gorm:"index"`
	// GORM relation
	User userModel `gorm:"foreignKey:UserEmail;references:Email;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name
func (postModel) TableName() string {
	return "posts"
}

// likeModel represents the database model for likes
type likeModel struct {
	UserEmail string `gorm:"primaryKey;column:user_email;not null"`
	PostID    string `gorm:"primaryKey;column:post_id;index;not null"`
	// GORM relations (using pointers to avoid circular reference issues)
	User *userModel `gorm:"foreignKey:UserEmail;references:Email;constraint:OnDelete:CASCADE"`
	Post *postModel `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name
func (likeModel) TableName() string {
	return "liked_posts"
}
