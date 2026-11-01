package repositories

import (
	"context"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
)

type AuthorizedPersonRepository interface {
	ListByStudentID(ctx context.Context, studentID int64) ([]entities.AuthorizedPerson, error)
	GetByID(ctx context.Context, id int64) (entities.AuthorizedPerson, error)
	Create(ctx context.Context, person entities.AuthorizedPerson) (entities.AuthorizedPerson, error)
	Update(ctx context.Context, person entities.AuthorizedPerson) (entities.AuthorizedPerson, error)
	Delete(ctx context.Context, id int64) error
}
