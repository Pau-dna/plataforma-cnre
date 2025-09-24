package postgres

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Get the underlying *sql.DB to configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Optimize connection pool settings for better performance
	sqlDB.SetMaxOpenConns(25)                 // Maximum connections in the pool
	sqlDB.SetMaxIdleConns(10)                 // Maximum idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime
	sqlDB.SetConnMaxIdleTime(2 * time.Minute) // Maximum idle time

	return db, nil
}
