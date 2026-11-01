package services

import (
	"context"
	"strings"

	"github.com/camilo-atb/aulaviva/backend/internal/application/dto"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/repositories"
)

type AuthorizedPersonService struct {
	students repositories.StudentRepository
	people   repositories.AuthorizedPersonRepository
}

func NewAuthorizedPersonService(students repositories.StudentRepository, people repositories.AuthorizedPersonRepository) *AuthorizedPersonService {
	return &AuthorizedPersonService{students: students, people: people}
}

func (s *AuthorizedPersonService) Create(ctx context.Context, studentID int64, input dto.AuthorizedPersonInput) (entities.AuthorizedPerson, error) {
	if _, err := s.students.GetByID(ctx, studentID); err != nil {
		return entities.AuthorizedPerson{}, err
	}
	person := toEntity(studentID, input)
	if !person.Validate() {
		return entities.AuthorizedPerson{}, errors.ErrInvalidInput
	}
	return s.people.Create(ctx, person)
}

func (s *AuthorizedPersonService) Update(ctx context.Context, id int64, input dto.AuthorizedPersonInput) (entities.AuthorizedPerson, error) {
	existing, err := s.people.GetByID(ctx, id)
	if err != nil {
		return entities.AuthorizedPerson{}, err
	}
	person := toEntity(existing.StudentID, input)
	person.ID = id
	if !person.Validate() {
		return entities.AuthorizedPerson{}, errors.ErrInvalidInput
	}
	return s.people.Update(ctx, person)
}

func (s *AuthorizedPersonService) Delete(ctx context.Context, id int64) error {
	return s.people.Delete(ctx, id)
}

func toEntity(studentID int64, input dto.AuthorizedPersonInput) entities.AuthorizedPerson {
	return entities.NormalizeAuthorizedPerson(entities.AuthorizedPerson{
		StudentID: studentID, FullName: strings.TrimSpace(input.FullName),
		DocumentType: input.DocumentType, DocumentNumber: strings.TrimSpace(input.DocumentNumber),
		Relationship: strings.TrimSpace(input.Relationship), Phone: strings.TrimSpace(input.Phone),
	})
}
