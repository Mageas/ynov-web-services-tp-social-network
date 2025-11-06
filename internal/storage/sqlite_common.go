package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteStores struct {
	DB *gorm.DB
}

// GORM models
type userModel struct {
	Email        string `gorm:"primaryKey"`
	SaltHex      string
	PasswordHash string
}

type postModel struct {
	ID        string `gorm:"primaryKey"`
	Author    string `gorm:"index"`
	Content   string
	CreatedAt int64
}

type likedPostModel struct {
	UserEmail string `gorm:"primaryKey;column:user_email"`
	PostID    string `gorm:"primaryKey;column:post_id"`
}

func (likedPostModel) TableName() string { return "liked_posts" }

func (userModel) TableName() string { return "users" }
func (postModel) TableName() string { return "posts" }

func OpenSQLite(path string) (*SQLiteStores, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	s := &SQLiteStores{DB: db}
	if err := s.migrate(); err != nil {
		if sqlDB, derr := db.DB(); derr == nil {
			_ = sqlDB.Close()
		}
		return nil, err
	}
	return s, nil
}

func (s *SQLiteStores) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *SQLiteStores) migrate() error {
	return s.DB.AutoMigrate(&userModel{}, &postModel{}, &likedPostModel{})
}
