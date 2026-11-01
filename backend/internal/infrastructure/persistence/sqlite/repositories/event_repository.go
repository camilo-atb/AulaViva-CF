package repositories

import (
	"context"
	"database/sql"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) CreateWithAuthorizations(ctx context.Context, event entities.Event, items []entities.StudentAuthorization) (entities.Event, []entities.Authorization, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return entities.Event{}, nil, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `INSERT INTO events (course_id, name, description, event_date) VALUES (?, ?, ?, ?)`,
		event.CourseID, event.Name, event.Description, event.EventDate)
	if err != nil {
		return entities.Event{}, nil, err
	}
	event.ID, err = result.LastInsertId()
	if err != nil {
		return entities.Event{}, nil, err
	}
	authorizations := make([]entities.Authorization, 0, len(items))
	for _, item := range items {
		result, err := tx.ExecContext(ctx, `INSERT INTO authorizations (event_id, student_id, request_token, status) VALUES (?, ?, ?, 'pending')`,
			event.ID, item.Student.ID, item.Authorization.RequestToken)
		if err != nil {
			return entities.Event{}, nil, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return entities.Event{}, nil, err
		}
		item.Authorization.ID = id
		item.Authorization.EventID = event.ID
		authorizations = append(authorizations, item.Authorization)
	}
	if err := tx.Commit(); err != nil {
		return entities.Event{}, nil, err
	}
	return event, authorizations, nil
}

func (r *EventRepository) List(ctx context.Context) ([]entities.Event, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT e.id, e.course_id, c.name, e.name, e.description, e.event_date
		FROM events e JOIN courses c ON c.id = e.course_id ORDER BY e.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []entities.Event
	for rows.Next() {
		var event entities.Event
		if err := rows.Scan(&event.ID, &event.CourseID, &event.CourseName, &event.Name, &event.Description, &event.EventDate); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (r *EventRepository) GetDetails(ctx context.Context, id int64) (entities.Event, []entities.Authorization, error) {
	var event entities.Event
	err := r.db.QueryRowContext(ctx, `SELECT e.id, e.course_id, c.name, e.name, e.description, e.event_date
		FROM events e JOIN courses c ON c.id = e.course_id WHERE e.id = ?`, id).
		Scan(&event.ID, &event.CourseID, &event.CourseName, &event.Name, &event.Description, &event.EventDate)
	if err == sql.ErrNoRows {
		return entities.Event{}, nil, domainerrors.ErrNotFound
	}
	if err != nil {
		return entities.Event{}, nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT a.id, a.event_id, a.student_id, s.full_name, a.request_token, a.status, a.responded_at
		FROM authorizations a JOIN students s ON s.id = a.student_id
		WHERE a.event_id = ? ORDER BY a.id`, id)
	if err != nil {
		return entities.Event{}, nil, err
	}
	defer rows.Close()
	var authorizations []entities.Authorization
	for rows.Next() {
		var item entities.Authorization
		var respondedAt sql.NullString
		if err := rows.Scan(&item.ID, &item.EventID, &item.StudentID, &item.StudentName, &item.RequestToken, &item.Status, &respondedAt); err != nil {
			return entities.Event{}, nil, err
		}
		if respondedAt.Valid {
			value, err := parseTime(respondedAt.String)
			if err != nil {
				return entities.Event{}, nil, err
			}
			item.RespondedAt = &value
		}
		authorizations = append(authorizations, item)
	}
	return event, authorizations, rows.Err()
}
