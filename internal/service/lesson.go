package service

import (
	"context"

	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/AdilzhanZh/LMS_backend/internal/repository"
)

type LessonService struct {
	repo repository.LessonRepo
}

func NewLessonService(repo repository.LessonRepo) *LessonService {
	return &LessonService{repo: repo}
}

func (ls *LessonService) GetAll(ctx context.Context) ([]models.Lesson, error) {
	return ls.repo.GetAll(ctx)
}

func (ls *LessonService) GetLessonByID(ctx context.Context, ID int) (models.Lesson, error) {
	return ls.repo.GetByID(ctx, ID)
}

func (ls *LessonService) DeleteByID(ctx context.Context, ID int) error {
	return ls.repo.DeleteByID(ctx, ID)
}

func (ls *LessonService) Create(ctx context.Context, lesson models.CreateLesson) (int, error) {
	return ls.repo.Create(ctx, lesson)
}

func (ls *LessonService) Update(ctx context.Context, id int, lesson models.UpdateLesson) (int, error) {
	return ls.repo.Update(ctx, id, lesson)
}
