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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type LessonRepo interface {
	GetAll(ctx context.Context) ([]models.Lesson, error)
	GetByID(ctx context.Context, ID int) (models.Lesson, error)
	DeleteByID(ctx context.Context, ID int) error
	Create(ctx context.Context, lesson models.CreateLesson) (int, error)
	Update(ctx context.Context, id int, lesson models.UpdateLesson) (int, error)
}

type PsgLessonRepo struct {
	db *sqlx.DB
}

func NewPsgLessonRepo(db *sqlx.DB) *PsgLessonRepo {
	return &PsgLessonRepo{
		db: db,
	}
}

func (plr *PsgLessonRepo) Create(ctx context.Context, lesson models.CreateLesson) (int, error) {
	query := `
		INSERT INTO lessons (
			course_id, title, content, video_url, duration, position,
			is_preview, created_at, updated_at
		) VALUES (
			:course_id, :title, :content, :video_url, :duration, :position,
			:is_preview, :created_at, :updated_at
		)
		RETURNING id
	`

	lesson.CreatedAt = utils.Now()
	lesson.UpdatedAt = utils.Now()

	stmt, err := plr.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("prepare query error: %w", err)
	}
	defer stmt.Close()

	var id int
	err = stmt.Get(&id, lesson)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return 0, models.ErrCourseNotFound
		}
		return 0, fmt.Errorf("create lesson error: %w", err)
	}

	return id, nil
}

func (plr *PsgLessonRepo) DeleteByID(ctx context.Context, ID int) error {
	query := `
		UPDATE lessons
		SET deleted_at = NOW(),
			updated_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`

	result, err := plr.db.ExecContext(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("delete lesson error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete lesson rowsAffected error: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrLessonNotFound
	}

	return nil
}

func (plr *PsgLessonRepo) GetByID(ctx context.Context, ID int) (models.Lesson, error) {
	var lesson models.Lesson

	query := `
		SELECT id, course_id, title, content, video_url, duration, position,
		is_preview, created_at, updated_at, deleted_at
		FROM lessons
		WHERE id = $1
		AND deleted_at IS NULL
		LIMIT 1
	`

	err := plr.db.GetContext(ctx, &lesson, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("Error no such id", "db", err)
			return models.Lesson{}, models.ErrLessonNotFound
		}
		slog.Error("Error with db", "db", err)
		return models.Lesson{}, fmt.Errorf("failed to get lesson by id err: %w", err)
	}

	return lesson, nil
}

func (plr *PsgLessonRepo) GetAll(ctx context.Context) ([]models.Lesson, error) {
	var lessons []models.Lesson

	query := `
		SELECT id, course_id, title, content, video_url, duration, position,
		is_preview, created_at, updated_at, deleted_at
		FROM lessons
		WHERE deleted_at IS NULL
		ORDER BY course_id ASC, position ASC, created_at DESC
	`

	err := plr.db.SelectContext(ctx, &lessons, query)
	if err != nil {
		slog.Error("failed to get lessons", "db", err)
		return nil, err
	}

	return lessons, nil
}

func (plr *PsgLessonRepo) Update(ctx context.Context, id int, lesson models.UpdateLesson) (int, error) {
	var setParts []string
	var args []interface{}
	argID := 1

	if lesson.CourseID != nil {
		setParts = append(setParts, fmt.Sprintf("course_id = $%d", argID))
		args = append(args, *lesson.CourseID)
		argID++
	}

	if lesson.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argID))
		args = append(args, *lesson.Title)
		argID++
	}

	if lesson.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content = $%d", argID))
		args = append(args, *lesson.Content)
		argID++
	}

	if lesson.VideoURL != nil {
		setParts = append(setParts, fmt.Sprintf("video_url = $%d", argID))
		args = append(args, *lesson.VideoURL)
		argID++
	}

	if lesson.Duration != nil {
		setParts = append(setParts, fmt.Sprintf("duration = $%d", argID))
		args = append(args, *lesson.Duration)
		argID++
	}

	if lesson.Position != nil {
		setParts = append(setParts, fmt.Sprintf("position = $%d", argID))
		args = append(args, *lesson.Position)
		argID++
	}

	if lesson.IsPreview != nil {
		setParts = append(setParts, fmt.Sprintf("is_preview = $%d", argID))
		args = append(args, *lesson.IsPreview)
		argID++
	}

	if len(setParts) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, utils.Now())
	argID++

	query := fmt.Sprintf(`
		UPDATE lessons
		SET %s
		WHERE id = $%d
		AND deleted_at IS NULL
		RETURNING id
	`, strings.Join(setParts, ", "), argID)

	args = append(args, id)

	var updatedID int
	if err := plr.db.GetContext(ctx, &updatedID, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrLessonNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return 0, models.ErrCourseNotFound
		}

		return 0, fmt.Errorf("update lesson: %w", err)
	}

	return updatedID, nil
}
