package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"agentic/commerce/config"
	"agentic/commerce/internal/infrastructure/database"
	"agentic/commerce/pkg/logger"

	"go.uber.org/fx"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

func registerDBShutdown(lc fx.Lifecycle, db *gorm.DB) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			fmt.Println("🛑 Gracefully stopping Database...")
			err := database.ShutdownGormDB(db)
			if err != nil {
				_ = fmt.Errorf("cant close database %v", err)
			}
			return nil
		},
	})
}

var DatabaseModule = fx.Module(
	"database",
	fx.Provide(NewDatabase),
	fx.Invoke(autoMigrate),
	fx.Invoke(registerDBShutdown),
)

type DB struct{}

func NewDatabase(cfg *config.DbConfig, appLogger *logger.AppLogger) (*gorm.DB, error) {
	lo := appLogger.WithScope(DB{})

	customLogger := NewZerologGormLogger(lo, gormLogger.LogLevel(lo.GetLogLevel()))

	db, err := database.NewGormPostgresConnection(cfg)

	if err != nil {
		lo.Fatal("error in new connection %v\n", err)
		return nil, err
	}

	db.Logger = customLogger

	return db, nil
}

func AutoMigrate(db *gorm.DB, models ...interface{}) {
	var err error

	nonConstraintGorm, err := gorm.Open(db.Dialector, &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(fmt.Errorf("cannot access db in auto-migrating %w", err))
	}
	defer func() {
		nonConstraintDB, _ := nonConstraintGorm.DB()
		_ = nonConstraintDB.Close()
	}()

	err = nonConstraintGorm.AutoMigrate(models...)
	if err != nil {
		panic(fmt.Errorf("error in non-constraint auto-migrating %w", err))
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		panic(fmt.Errorf("error in auto-migrating %w", err))
	}
}

// AutoMigrate performs database migration for all registered models.
func autoMigrate(db *gorm.DB, registry database.EntityRegistry) {

	db = db.Debug()

	// Convert models to []interface{} for GORM
	models := make([]interface{}, len(registry.Models))
	for i, m := range registry.Models {
		models[i] = m
	}

	// Execute migration
	AutoMigrate(db, models...)

	log.Println("Auto migration completed successfully")
}

type ZerologGormLogger struct {
	log        *logger.AppLogger
	level      gormLogger.LogLevel
	slowThresh time.Duration
}

func NewZerologGormLogger(appLogger *logger.AppLogger, level gormLogger.LogLevel) *ZerologGormLogger {
	return &ZerologGormLogger{
		log:        appLogger,
		level:      level,
		slowThresh: 200 * time.Millisecond,
	}
}

func (l *ZerologGormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.level = level
	return l
}

func (l *ZerologGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormLogger.Info {
		l.log.Info(msg, "args", data)
	}
}

func (l *ZerologGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormLogger.Warn {
		l.log.Warn(msg, "args", data)
	}
}

func (l *ZerologGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormLogger.Error {
		l.log.Error(msg, nil)
	}
}

func (l *ZerologGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	e := l.log.With("file", utils.FileWithLineNum()).
		With("sql", sql).
		With("rows", rows).
		With("elapsed", elapsed)

	switch {
	case err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound):
		e.Error("gorm query error", err)
	case elapsed > l.slowThresh && l.level >= gormLogger.Warn:
		e.Warn("slow query")
	case l.level >= gormLogger.Info:
		e.Debug("gorm query")
	}
}
