package models

import "errors"

var ErrSlugAlreadyExists = errors.New("slug already exists error")
var ErrTeacherNotFound = errors.New("teacher not found error")
var ErrCourseNotFound = errors.New("course not found error")
var ErrLessonNotFound = errors.New("lesson not found error")
