package dto

import "github.com/camilo-atb/aulaviva/backend/internal/domain/entities"

type CreateEventInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	EventDate   string `json:"eventDate"`
	CourseID    int64  `json:"courseId"`
}

type AuthorizationSummary struct {
	Pending       int `json:"pending"`
	Authorized    int `json:"authorized"`
	NotAuthorized int `json:"notAuthorized"`
}

type EventResponse struct {
	entities.Event
	Authorizations []entities.Authorization `json:"authorizations,omitempty"`
	Summary        AuthorizationSummary     `json:"summary,omitempty"`
}

type RequestResponse struct {
	Authorization entities.Authorization `json:"authorization"`
	Event         entities.Event         `json:"event"`
	Student       entities.Student       `json:"student"`
}

type SubmitResponseInput struct {
	Status entities.AuthorizationStatus `json:"status"`
}

type AuthorizedPersonInput struct {
	FullName       string                `json:"fullName"`
	DocumentType   entities.DocumentType `json:"documentType"`
	DocumentNumber string                `json:"documentNumber"`
	Relationship   string                `json:"relationship"`
	Phone          string                `json:"phone"`
}

type StudentResponse struct {
	entities.Student
	AuthorizedPeople []entities.AuthorizedPerson `json:"authorizedPeople"`
}
