package entities

import (
	"time"

	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
)

type AuthorizationStatus string

const (
	StatusPending       AuthorizationStatus = "pending"
	StatusAuthorized    AuthorizationStatus = "authorized"
	StatusNotAuthorized AuthorizationStatus = "not_authorized"
)

type Authorization struct {
	ID           int64               `json:"id"`
	EventID      int64               `json:"eventId"`
	StudentID    int64               `json:"studentId"`
	StudentName  string              `json:"studentName,omitempty"`
	RequestToken string              `json:"requestToken"`
	Status       AuthorizationStatus `json:"status"`
	RespondedAt  *time.Time          `json:"respondedAt,omitempty"`
}

func (a Authorization) ValidateResponse(status AuthorizationStatus) error {
	if status != StatusAuthorized && status != StatusNotAuthorized {
		return domainerrors.ErrInvalidStatus
	}
	return nil
}
