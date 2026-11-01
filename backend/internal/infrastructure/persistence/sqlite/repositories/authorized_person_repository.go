package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
)

type AuthorizedPersonRepository struct {
	db *sql.DB
}

func NewAuthorizedPersonRepository(db *sql.DB) *AuthorizedPersonRepository {
	return &AuthorizedPersonRepository{db: db}
}

func (r *AuthorizedPersonRepository) ListByStudentID(ctx context.Context, studentID int64) ([]entities.AuthorizedPerson, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, student_id, full_name, document_type, document_number, relationship, phone
		FROM authorized_people WHERE student_id = ? ORDER BY full_name, id`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	people := make([]entities.AuthorizedPerson, 0)
	for rows.Next() {
		person, err := scanAuthorizedPerson(rows)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}
	return people, rows.Err()
}

func (r *AuthorizedPersonRepository) GetByID(ctx context.Context, id int64) (entities.AuthorizedPerson, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, student_id, full_name, document_type, document_number, relationship, phone
		FROM authorized_people WHERE id = ?`, id)
	person, err := scanAuthorizedPerson(row)
	if err == sql.ErrNoRows {
		return entities.AuthorizedPerson{}, domainerrors.ErrNotFound
	}
	return person, err
}

func (r *AuthorizedPersonRepository) Create(ctx context.Context, person entities.AuthorizedPerson) (entities.AuthorizedPerson, error) {
	result, err := r.db.ExecContext(ctx, `INSERT INTO authorized_people
		(student_id, full_name, document_type, document_number, relationship, phone)
		VALUES (?, ?, ?, ?, ?, ?)`, person.StudentID, person.FullName, person.DocumentType,
		person.DocumentNumber, person.Relationship, person.Phone)
	if isConstraintError(err) {
		return entities.AuthorizedPerson{}, domainerrors.ErrDuplicate
	}
	if err != nil {
		return entities.AuthorizedPerson{}, err
	}
	person.ID, err = result.LastInsertId()
	return person, err
}

func (r *AuthorizedPersonRepository) Update(ctx context.Context, person entities.AuthorizedPerson) (entities.AuthorizedPerson, error) {
	result, err := r.db.ExecContext(ctx, `UPDATE authorized_people SET full_name = ?, document_type = ?,
		document_number = ?, relationship = ?, phone = ? WHERE id = ?`, person.FullName, person.DocumentType,
		person.DocumentNumber, person.Relationship, person.Phone, person.ID)
	if isConstraintError(err) {
		return entities.AuthorizedPerson{}, domainerrors.ErrDuplicate
	}
	if err != nil {
		return entities.AuthorizedPerson{}, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return entities.AuthorizedPerson{}, err
	}
	if affected == 0 {
		return entities.AuthorizedPerson{}, domainerrors.ErrNotFound
	}
	return person, nil
}

func (r *AuthorizedPersonRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM authorized_people WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domainerrors.ErrNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanAuthorizedPerson(row scanner) (entities.AuthorizedPerson, error) {
	var person entities.AuthorizedPerson
	err := row.Scan(&person.ID, &person.StudentID, &person.FullName, &person.DocumentType,
		&person.DocumentNumber, &person.Relationship, &person.Phone)
	return person, err
}

func isConstraintError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToUpper(err.Error())
	return strings.Contains(message, "UNIQUE") || strings.Contains(message, "CHECK") || strings.Contains(message, "FOREIGN KEY")
}
