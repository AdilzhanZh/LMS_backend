package service

import (
	"context"
	"errors"

	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/AdilzhanZh/LMS_backend/internal/repository"
)

type LessonService struct {
	repo       repository.LessonRepo
	courseRepo repository.CourseRepo
}

func NewLessonService(repo repository.LessonRepo, courseRepo repository.CourseRepo) *LessonService {
	return &LessonService{
		repo:       repo,
		courseRepo: courseRepo,
	}
}

func (ls *LessonService) GetAll(ctx context.Context) ([]models.Lesson, error) {
	return ls.repo.GetAll(ctx)
}

func (ls *LessonService) GetLessonByID(ctx context.Context, ID int) (models.Lesson, error) {
	return ls.repo.GetByID(ctx, ID)
}

func (ls *LessonService) DeleteByID(ctx context.Context, ID int) error {

	lesson, err := ls.repo.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	course, err := ls.courseRepo.GetByID(ctx, lesson.CourseID)
	if err != nil {
		return err
	}

	if course.IsActive {
		return errors.New("cannot delete lesson inside active course")
	}

	return ls.repo.DeleteByID(ctx, ID)
}

func (ls *LessonService) Create(ctx context.Context, lesson models.CreateLesson) (int, error) {

	_, err := ls.courseRepo.GetByID(ctx, lesson.CourseID)
	if err != nil {
		return 0, err
	}

	return ls.repo.Create(ctx, lesson)
}

func (ls *LessonService) Update(ctx context.Context, id int, lesson models.UpdateLesson) (int, error) {
	return ls.repo.Update(ctx, id, lesson)
}
