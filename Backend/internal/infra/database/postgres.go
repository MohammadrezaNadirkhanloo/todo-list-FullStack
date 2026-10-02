package database

import (
	"context"
	"fmt"
	"time"

	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func Connect(ctx context.Context, cfg config.PostgresConfig) (*DB, error) {
	gormCfg := &gorm.Config{
		// Logger: NewGormLogger(log, cfg.SlowQueryThreshold),
		NamingStrategy:                           nil,
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
	}

	gdb, err := gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("database: failed to connect to PostgreSQL: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("database: failed to get underlying connection: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("database: پاسخ ندادن PostgreSQL: %w", err)
	}

	// log.Info(ctx, "اتصال به PostgreSQL برقرار شد",
	// 	logging.Cat(logging.CategoryPostgres),
	// 	logging.Sub(logging.SubConnection),
	// 	logging.F("host", cfg.Host),
	// 	logging.F("database", cfg.DBName),
	// 	logging.F("max_open_conns", cfg.MaxOpenConns),
	// )

	fmt.Println("connecting Database")

	return &DB{DB: gdb}, nil
}

func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *DB) HealthCheck(ctx context.Context) error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
