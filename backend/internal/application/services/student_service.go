package services

import (
	"context"

	"github.com/camilo-atb/aulaviva/backend/internal/application/dto"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/repositories"
)

type StudentService struct {
	students repositories.StudentRepository
	people   repositories.AuthorizedPersonRepository
}

func NewStudentService(students repositories.StudentRepository, people repositories.AuthorizedPersonRepository) *StudentService {
	return &StudentService{students: students, people: people}
}

func (s *StudentService) List(ctx context.Context, search string) ([]dto.StudentResponse, error) {
	students, err := s.students.List(ctx, search)
	if err != nil {
		return nil, err
	}
	result := make([]dto.StudentResponse, 0, len(students))
	for _, student := range students {
		people, err := s.people.ListByStudentID(ctx, student.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, studentResponse(student, people))
	}
	return result, nil
}

func (s *StudentService) Get(ctx context.Context, id int64) (dto.StudentResponse, error) {
	student, err := s.students.GetByID(ctx, id)
	if err != nil {
		return dto.StudentResponse{}, err
	}
	people, err := s.people.ListByStudentID(ctx, id)
	if err != nil {
		return dto.StudentResponse{}, err
	}
	return studentResponse(student, people), nil
}

func (s *StudentService) ListPeople(ctx context.Context, studentID int64) ([]entities.AuthorizedPerson, error) {
	if _, err := s.students.GetByID(ctx, studentID); err != nil {
		return nil, err
	}
	return s.people.ListByStudentID(ctx, studentID)
}

func studentResponse(student entities.Student, people []entities.AuthorizedPerson) dto.StudentResponse {
	if people == nil {
		people = []entities.AuthorizedPerson{}
	}
	return dto.StudentResponse{Student: student, AuthorizedPeople: people}
}
