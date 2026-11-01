package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
)

type AuthorizationRepository struct {
	db *sql.DB
}

func NewAuthorizationRepository(db *sql.DB) *AuthorizationRepository {
	return &AuthorizationRepository{db: db}
}

func (r *AuthorizationRepository) GetByToken(ctx context.Context, token string) (entities.Authorization, entities.Event, entities.Student, error) {
	var authorization entities.Authorization
	var event entities.Event
	var student entities.Student
	var respondedAt sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT a.id, a.event_id, a.student_id, a.request_token, a.status, a.responded_at,
		e.course_id, e.name, e.description, e.event_date, s.course_id, s.full_name
		FROM authorizations a JOIN events e ON e.id = a.event_id JOIN students s ON s.id = a.student_id
		WHERE a.request_token = ?`, token).Scan(
		&authorization.ID, &authorization.EventID, &authorization.StudentID, &authorization.RequestToken,
		&authorization.Status, &respondedAt, &event.CourseID, &event.Name, &event.Description,
		&event.EventDate, &student.CourseID, &student.FullName)
	if err == sql.ErrNoRows {
		return entities.Authorization{}, entities.Event{}, entities.Student{}, domainerrors.ErrNotFound
	}
	if err != nil {
		return entities.Authorization{}, entities.Event{}, entities.Student{}, err
	}
	event.ID = authorization.EventID
	student.ID = authorization.StudentID
	if respondedAt.Valid {
		value, err := parseTime(respondedAt.String)
		if err != nil {
			return entities.Authorization{}, entities.Event{}, entities.Student{}, err
		}
		authorization.RespondedAt = &value
	}
	return authorization, event, student, nil
}

func (r *AuthorizationRepository) Respond(ctx context.Context, token string, status entities.AuthorizationStatus, respondedAt time.Time) error {
	result, err := r.db.ExecContext(ctx, `UPDATE authorizations SET status = ?, responded_at = ?
		WHERE request_token = ? AND status = 'pending'`, status, respondedAt.Format(time.RFC3339Nano), token)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 1 {
		return nil
	}
	var exists int
	err = r.db.QueryRowContext(ctx, `SELECT 1 FROM authorizations WHERE request_token = ?`, token).Scan(&exists)
	if err == sql.ErrNoRows {
		return domainerrors.ErrNotFound
	}
	if err != nil {
		return err
	}
	return domainerrors.ErrAlreadyAnswered
}

func parseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, value)
}
