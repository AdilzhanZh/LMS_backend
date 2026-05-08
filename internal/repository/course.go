package repository

import (
	"context"

	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type CourseRepo interface {
	GetAll(ctx context.Context) ([]models.Course, error)
	//TODO implement other methods
}

type PsgCourseRepo struct {
	db *sqlx.DB
}

func NewPsgCourseRepo(db *sqlx.DB) *PsgCourseRepo {
	return &PsgCourseRepo{
		db: db,
	}
}

func (pcr *PsgCourseRepo) GetAll(ctx context.Context) ([]models.Course, error) {
	var courses []models.Course

	query := `
		SELECT id, title, description, slug, price, duration, level,
		is_active, teacher_id, created_at, updated_at, deleted_at
		FROM courses
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	err := pcr.db.SelectContext(ctx, &courses, query)
	if err != nil {
		return nil, err
	}

	return courses, nil
}
