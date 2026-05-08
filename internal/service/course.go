package service

import (
	"context"

	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/AdilzhanZh/LMS_backend/internal/repository"
)

type CourseService struct {
	repo repository.CourseRepo // interface
}

func NewCourseService(repo repository.CourseRepo) *CourseService {
	return &CourseService{repo: repo}
}

func (cs *CourseService) GetAll(ctx context.Context) ([]models.Course, error) {
	return cs.repo.GetAll(ctx)
}

func (cs *CourseService) GetCourseByID(ctx context.Context, ID int) (models.Course, error) {
	return cs.repo.GetByID(ctx, ID)
}

func (cs *CourseService) DeleteByID(ctx context.Context, ID int) error {
	return cs.repo.DeleteByID(ctx, ID)
}

func (cs *CourseService) Create(ctx context.Context, course models.CreateCourse) (int, error) {
	return cs.repo.Create(ctx, course)
}
