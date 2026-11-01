package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/camilo-atb/aulaviva/backend/internal/application/dto"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	"github.com/camilo-atb/aulaviva/backend/internal/domain/repositories"
)

type EventService struct {
	courses repositories.CourseRepository
	events  repositories.EventRepository
}

func NewEventService(courses repositories.CourseRepository, events repositories.EventRepository) *EventService {
	return &EventService{courses: courses, events: events}
}

func (s *EventService) Create(ctx context.Context, input dto.CreateEventInput) (dto.EventResponse, error) {
	event := entities.Event{
		Name: strings.TrimSpace(input.Name), Description: strings.TrimSpace(input.Description),
		EventDate: input.EventDate, CourseID: input.CourseID,
	}
	if err := event.Validate(); err != nil {
		return dto.EventResponse{}, err
	}
	course, err := s.courses.GetByID(ctx, event.CourseID)
	if err != nil {
		return dto.EventResponse{}, err
	}
	students, err := s.courses.ListStudents(ctx, course.ID)
	if err != nil {
		return dto.EventResponse{}, err
	}
	event.CourseName = course.Name
	authorizations := make([]entities.Authorization, 0, len(students))
	for _, student := range students {
		token, err := newToken()
		if err != nil {
			return dto.EventResponse{}, err
		}
		authorizations = append(authorizations, entities.Authorization{
			StudentID: student.ID, RequestToken: token, Status: entities.StatusPending,
		})
	}
	created, saved, err := s.events.CreateWithAuthorizations(ctx, event, authorizationsToStudents(authorizations, students))
	if err != nil {
		return dto.EventResponse{}, err
	}
	return dto.EventResponse{Event: created, Authorizations: saved, Summary: summarize(saved)}, nil
}

func (s *EventService) List(ctx context.Context) ([]dto.EventResponse, error) {
	events, err := s.events.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.EventResponse, 0, len(events))
	for _, event := range events {
		result = append(result, dto.EventResponse{Event: event})
	}
	return result, nil
}

func (s *EventService) Get(ctx context.Context, id int64) (dto.EventResponse, error) {
	event, authorizations, err := s.events.GetDetails(ctx, id)
	if err != nil {
		return dto.EventResponse{}, err
	}
	return dto.EventResponse{Event: event, Authorizations: authorizations, Summary: summarize(authorizations)}, nil
}

func summarize(items []entities.Authorization) dto.AuthorizationSummary {
	var summary dto.AuthorizationSummary
	for _, item := range items {
		switch item.Status {
		case entities.StatusPending:
			summary.Pending++
		case entities.StatusAuthorized:
			summary.Authorized++
		case entities.StatusNotAuthorized:
			summary.NotAuthorized++
		}
	}
	return summary
}

func authorizationsToStudents(items []entities.Authorization, students []entities.Student) []entities.StudentAuthorization {
	result := make([]entities.StudentAuthorization, len(items))
	for i := range items {
		result[i] = entities.StudentAuthorization{Authorization: items[i], Student: students[i]}
	}
	return result
}

func newToken() (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}
