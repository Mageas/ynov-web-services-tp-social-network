package sqlite

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
type DB struct {
	conn *gorm.DB
}

// New creates a new database connection
func New(path string) (*DB, error) {
	conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// migrate runs database migrations
func (db *DB) migrate() error {
	return db.conn.AutoMigrate(
		&userModel{},
		&postModel{},
		&likeModel{},
	)
}

// GetConn returns the underlying GORM connection
func (db *DB) GetConn() *gorm.DB {
	return db.conn
}
