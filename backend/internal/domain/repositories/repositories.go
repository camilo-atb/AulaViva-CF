package repositories

import (
	"context"
	"time"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
)

type CourseRepository interface {
	List(ctx context.Context) ([]entities.Course, error)
	GetByID(ctx context.Context, id int64) (entities.Course, error)
	ListStudents(ctx context.Context, courseID int64) ([]entities.Student, error)
}

type EventRepository interface {
	CreateWithAuthorizations(ctx context.Context, event entities.Event, authorizations []entities.StudentAuthorization) (entities.Event, []entities.Authorization, error)
	List(ctx context.Context) ([]entities.Event, error)
	GetDetails(ctx context.Context, id int64) (entities.Event, []entities.Authorization, error)
}

type AuthorizationRepository interface {
	GetByToken(ctx context.Context, token string) (entities.Authorization, entities.Event, entities.Student, error)
	Respond(ctx context.Context, token string, status entities.AuthorizationStatus, respondedAt time.Time) error
}
