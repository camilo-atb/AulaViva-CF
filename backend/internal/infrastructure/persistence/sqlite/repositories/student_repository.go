package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) List(ctx context.Context, search string) ([]entities.Student, error) {
	query := `SELECT id, course_id, full_name FROM students`
	args := []any{}
	if strings.TrimSpace(search) != "" {
		query += ` WHERE LOWER(full_name) LIKE LOWER(?)`
		args = append(args, "%"+strings.TrimSpace(search)+"%")
	}
	query += ` ORDER BY full_name, id`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	students := make([]entities.Student, 0)
	for rows.Next() {
		var student entities.Student
		if err := rows.Scan(&student.ID, &student.CourseID, &student.FullName); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, rows.Err()
}

func (r *StudentRepository) GetByID(ctx context.Context, id int64) (entities.Student, error) {
	var student entities.Student
	err := r.db.QueryRowContext(ctx, `SELECT id, course_id, full_name FROM students WHERE id = ?`, id).
		Scan(&student.ID, &student.CourseID, &student.FullName)
	if err == sql.ErrNoRows {
		return entities.Student{}, domainerrors.ErrNotFound
	}
	return student, err
}
