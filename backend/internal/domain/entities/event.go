package entities

import (
	"strings"
	"time"

	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
)

type Event struct {
	ID          int64  `json:"id"`
	CourseID    int64  `json:"courseId"`
	CourseName  string `json:"courseName,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	EventDate   string `json:"eventDate"`
}

func (e Event) Validate() error {
	if strings.TrimSpace(e.Name) == "" || e.CourseID <= 0 {
		return domainerrors.ErrInvalidInput
	}
	if _, err := time.Parse("2006-01-02", e.EventDate); err != nil {
		return domainerrors.ErrInvalidInput
	}
	return nil
}
