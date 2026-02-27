package repository

import "github.com/AdilzhanZh/LMS_backend/internal/config"

type DB struct {
	connectionPath string
}

func NewPostgresDB(cfg *config.Config) (*DB, error) {
	db := &DB{}
	return db, nil
}
