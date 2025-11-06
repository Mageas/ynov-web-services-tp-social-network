package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// GORM-backed UserStore
type SQLiteUserStore struct{ db *gorm.DB }

func (s *SQLiteStores) Users() *SQLiteUserStore { return &SQLiteUserStore{db: s.DB} }

func (u *SQLiteUserStore) CreateUser(email, password string) error {
	var existing userModel
	if err := u.db.Where("LOWER(email) = LOWER(?)", email).First(&existing).Error; err == nil {
		return errors.New("user already exists")
	}
	// keep deterministic salt scheme compatible with previous version
	saltHex := hex.EncodeToString([]byte(email))
	salt, _ := hex.DecodeString(saltHex)
	sum := sha256.Sum256(append([]byte(password), salt...))
	rec := userModel{Email: strings.ToLower(email), SaltHex: saltHex, PasswordHash: hex.EncodeToString(sum[:])}
	return u.db.Create(&rec).Error
}

func (u *SQLiteUserStore) VerifyCredentials(email, password string) bool {
	var rec userModel
	if err := u.db.Where("LOWER(email) = LOWER(?)", email).First(&rec).Error; err != nil {
		return false
	}
	salt, err := hex.DecodeString(rec.SaltHex)
	if err != nil {
		return false
	}
	sum := sha256.Sum256(append([]byte(password), salt...))
	return hex.EncodeToString(sum[:]) == rec.PasswordHash
}
