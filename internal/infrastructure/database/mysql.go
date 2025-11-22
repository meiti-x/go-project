package database

import (
	"fmt"
	"time"

	"agentic/commerce/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewGormPostgresConnection(cfg *config.DbConfig) (*gorm.DB, error) {

	var db *gorm.DB
	var err error
	counter := 0
	var pid int

	schema.RegisterSerializer("enum", orsiniumEnumSerializer{})

	baseDSN := cfg.PostgresDSN()

	db, err = gorm.Open(postgres.Open(baseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot open database %s: %w", baseDSN, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("cannot get sql database %s: %w", baseDSN, err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConn)
	sqlDB.SetMaxIdleConns(cfg.IdleConn)
	sqlDB.SetConnMaxLifetime(cfg.Timeout)

	// Retry loop (like your MySQL version)
	for {
		<-time.NewTicker(cfg.DialTimeout).C
		counter++

		err = sqlDB.QueryRow("SELECT pg_backend_pid()").Scan(&pid)
		if err == nil {
			var version string
			err = sqlDB.QueryRow("SHOW server_version").Scan(&version)
			fmt.Println(fmt.Sprintf("PostgreSQL version: %s", version), zap.Error(err))
			break
		}

		fmt.Println(fmt.Sprintf("Cannot connect to PostgreSQL %s", baseDSN), zap.Error(err))
		if counter >= cfg.DialRetry {
			return nil, fmt.Errorf("cannot connect to PostgreSQL %s after %d retries", baseDSN, counter)
		}
	}

	fmt.Printf("Connected to %s PostgreSQL database: %s\n", cfg.Database, baseDSN)
	return db, nil
}
