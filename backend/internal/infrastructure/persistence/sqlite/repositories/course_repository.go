package repositories

import (
	"context"
	"database/sql"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
)

type CourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) List(ctx context.Context) ([]entities.Course, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM courses ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []entities.Course
	for rows.Next() {
		var course entities.Course
		if err := rows.Scan(&course.ID, &course.Name); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, rows.Err()
}

func (r *CourseRepository) GetByID(ctx context.Context, id int64) (entities.Course, error) {
	var course entities.Course
	err := r.db.QueryRowContext(ctx, `SELECT id, name FROM courses WHERE id = ?`, id).Scan(&course.ID, &course.Name)
	if err == sql.ErrNoRows {
		return entities.Course{}, domainerrors.ErrNotFound
	}
	return course, err
}

func (r *CourseRepository) ListStudents(ctx context.Context, courseID int64) ([]entities.Student, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, course_id, full_name FROM students WHERE course_id = ? ORDER BY id`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var students []entities.Student
	for rows.Next() {
		var student entities.Student
		if err := rows.Scan(&student.ID, &student.CourseID, &student.FullName); err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, rows.Err()
}
