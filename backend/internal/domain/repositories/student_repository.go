package repositories

import (
	"context"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
)

type StudentRepository interface {
	List(ctx context.Context, search string) ([]entities.Student, error)
	GetByID(ctx context.Context, id int64) (entities.Student, error)
}
