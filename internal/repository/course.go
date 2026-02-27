package repository

import "github.com/AdilzhanZh/LMS_backend/internal/models"

type CourseRepo interface {
	GetAll() ([]models.Course, error)
	//TODO implement other methods
}

type PsgCourseRepo struct {
	db *DB
}

func NewPsqCourseRepo(db *DB) *PsgCourseRepo {
	return &PsgCourseRepo{
		db: db,
	}
}

func (pcr *PsgCourseRepo) GetAll() ([]models.Course, error) {
	return []models.Course{
		{ID: 1, Name: "Go Basics"},
		{ID: 2, Name: "Nodejs Basics"},
		{ID: 3, Name: "Java Basics"},
		{ID: 4, Name: "Python Basics"},
	}, nil
}
