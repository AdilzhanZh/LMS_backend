package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AdilzhanZh/LMS_backend/internal/config"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	connectionPath string
}

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	slog.Info("PostgresSQL connected successfully")

	return db, nil

}
