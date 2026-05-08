package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/AdilzhanZh/LMS_backend/internal/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type CourseRepo interface {
	GetAll(ctx context.Context) ([]models.Course, error)
	GetByID(ctx context.Context, ID int) (models.Course, error)
	DeleteByID(ctx context.Context, ID int) error
	Create(ctx context.Context, course models.CreateCourse) (int, error)
	Update(ctx context.Context, id int, course models.UpdateCourse) (int, error)
}

type PsgCourseRepo struct {
	db *sqlx.DB
}

func NewPsgCourseRepo(db *sqlx.DB) *PsgCourseRepo {
	return &PsgCourseRepo{
		db: db,
	}
}

func (pcr *PsgCourseRepo) Create(ctx context.Context, course models.CreateCourse) (int, error) {
	query := `
		INSERT INTO courses (
			title, description, slug, price, duration, level, 
		    is_active, teacher_id, created_at, updated_at
		) VALUES (
		    :title, :description, :slug, :price, :duration, :level, 
		    :is_active, :teacher_id, :created_at, :updated_at
		)
		RETURNING id
	`
	course.CreatedAt = utils.Now()
	course.UpdatedAt = utils.Now()

	stmt, err := pcr.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare query error: %w", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.Get(&id, course)
	if err != nil {
		return 0, fmt.Errorf("create courses error: %w", err)
	}
	return id, nil
}

func (pcr *PsgCourseRepo) DeleteByID(ctx context.Context, ID int) error {
	query := `
		UPDATE courses
		SET deleted_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`

	result, err := pcr.db.ExecContext(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("delete course error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete course rowsAffected error: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (pcr *PsgCourseRepo) GetByID(ctx context.Context, ID int) (models.Course, error) {
	var course models.Course

	query := `
		SELECT id, title, description, slug, price, duration, level,
		is_active, teacher_id, created_at, updated_at, deleted_at
		FROM courses
		WHERE id = $1 
		AND deleted_at IS NULL
		LIMIT 1
	`

	err := pcr.db.GetContext(ctx, &course, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("Error no such id", "db", err)
			return models.Course{}, models.ErrNotFound
		}
		slog.Error("Error with db", "db", err)
		return models.Course{}, fmt.Errorf("failed to get course by id err: %w", err)
	}
	return course, nil
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
		slog.Error("failed to get courses", "db", err)
		return nil, err
	}

	return courses, nil
}

func (pcr *PsgCourseRepo) Update(ctx context.Context, id int, course models.UpdateCourse) (int, error) {
	var setParts []string  // "title", "price"
	var args []interface{} // "Golang", 10000
	argID := 1

	if course.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argID))
		args = append(args, *course.Title)
		argID++
	}

	if course.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argID))
		args = append(args, *course.Description)
		argID++
	}

	if course.Slug != nil {
		setParts = append(setParts, fmt.Sprintf("slug = $%d", argID))
		args = append(args, *course.Slug)
		argID++
	}

	if course.Price != nil {
		setParts = append(setParts, fmt.Sprintf("price = $%d", argID))
		args = append(args, *course.Price)
		argID++
	}

	if course.Duration != nil {
		setParts = append(setParts, fmt.Sprintf("duration = $%d", argID))
		args = append(args, *course.Duration)
		argID++
	}

	if course.Level != nil {
		setParts = append(setParts, fmt.Sprintf("level = $%d", argID))
		args = append(args, *course.Level)
		argID++
	}

	if course.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argID))
		args = append(args, *course.IsActive)
		argID++
	}

	if course.TeacherID != nil {
		setParts = append(setParts, fmt.Sprintf("teacher_id = $%d", argID))
		args = append(args, *course.TeacherID)
		argID++
	}

	if len(setParts) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, utils.Now())
	argID++

	//noinspection
	query := fmt.Sprintf(`
		UPDATE courses
		SET %s
		WHERE id = $%d
		AND deleted_at IS NULL
		RETURNING id
	`, strings.Join(setParts, ", "), argID)

	args = append(args, id)

	var updatedID int
	if err := pcr.db.GetContext(ctx, &updatedID, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNotFound
		}
		return 0, fmt.Errorf("update course: %w", err)
	}
	return updatedID, nil
}
