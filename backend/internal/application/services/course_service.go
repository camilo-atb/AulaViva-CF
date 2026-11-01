package services

import (
	"context"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/repositories"
)

type CourseService struct {
	courses repositories.CourseRepository
}

func NewCourseService(courses repositories.CourseRepository) *CourseService {
	return &CourseService{courses: courses}
}

func (s *CourseService) List(ctx context.Context) ([]entities.Course, error) {
	return s.courses.List(ctx)
}
