package service

import (
	"github.com/AdilzhanZh/LMS_backend/internal/models"
	"github.com/AdilzhanZh/LMS_backend/internal/repository"
)

type CourseService struct {
	repo repository.CourseRepo // interface
}

func NewCourseService(repo repository.CourseRepo) *CourseService {
	return &CourseService{repo: repo}
}

func (cs *CourseService) GetAll() ([]models.Course, error) {
	return cs.repo.GetAll()
}
