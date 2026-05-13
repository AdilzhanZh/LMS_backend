package service

import (
	"context"
	"errors"

	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/AdilzhanZh/LMS_backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type CourseService struct {
	repo       repository.CourseRepo // interface
	lessonRepo repository.LessonRepo
	enrollRepo repository.EnrollmentRepo
	db         *sqlx.DB
}

func NewCourseService(
	repo repository.CourseRepo,
	lessonRepo repository.LessonRepo,
	enrollRepo repository.EnrollmentRepo,
	db *sqlx.DB,
) *CourseService {
	return &CourseService{
		repo:       repo,
		lessonRepo: lessonRepo,
		enrollRepo: enrollRepo,
		db:         db,
	}
}

func (cs *CourseService) GetAll(ctx context.Context) ([]models.Course, error) {
	return cs.repo.GetAll(ctx)
}

func (cs *CourseService) GetCourseByID(ctx context.Context, ID int) (models.Course, error) {
	return cs.repo.GetByID(ctx, ID)
}

func (cs *CourseService) DeleteByID(ctx context.Context, id int) error {
	tx, err := cs.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err = cs.lessonRepo.DeleteByCourseIDTx(ctx, tx, id); err != nil {
		return err
	}

	if err = cs.enrollRepo.DeleteByCourseIDTx(ctx, tx, id); err != nil {
		return err
	}

	if err = cs.repo.DeleteByIDTx(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (cs *CourseService) Create(ctx context.Context, course models.CreateCourse) (int, error) {
	if err := course.Validate(); err != nil {
		return 0, err
	}

	course.IsActive = false // ishinde lesson bolmaiynsha false

	return cs.repo.Create(ctx, course)
}

func (cs *CourseService) Update(ctx context.Context, id int, course models.UpdateCourse) (int, error) {
	if err := course.Validate(); err != nil {
		return 0, err
	}

	if course.IsActive != nil && *course.IsActive == true {
		lessons, _ := cs.lessonRepo.GetByCourseID(ctx, id)
		if len(lessons) == 0 {
			return 0, errors.New("cannot activate course without lessons")
		}
	}

	return cs.repo.Update(ctx, id, course)
}
