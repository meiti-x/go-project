package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type GormDB func(context.Context) *gorm.DB

func CreateGormDB(db *gorm.DB) GormDB {
	return func(ctx context.Context) *gorm.DB {
		if db == nil {
			return nil
		}

		return db.WithContext(ctx)
	}
}

func ShutdownGormDB(db *gorm.DB) error {
	database, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	if err := database.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %w", err)
	}

	return nil
}
