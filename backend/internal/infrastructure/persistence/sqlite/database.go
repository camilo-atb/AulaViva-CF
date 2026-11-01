package sqlite

import (
	"context"
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed migrations/001_initial.sql
var initialMigration string

func Open(ctx context.Context, path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		return nil, err
	}
	if err := migrate(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, initialMigration)
	return err
}

func Seed(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM courses`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err := db.ExecContext(ctx, `
INSERT INTO courses (name) VALUES ('Primero A'), ('Segundo A');
INSERT INTO students (course_id, full_name)
SELECT id, 'Sofía Martínez' FROM courses WHERE name = 'Primero A';
INSERT INTO students (course_id, full_name)
SELECT id, 'Mateo Rodríguez' FROM courses WHERE name = 'Primero A';
INSERT INTO students (course_id, full_name)
SELECT id, 'Valentina Gómez' FROM courses WHERE name = 'Segundo A';
INSERT INTO students (course_id, full_name)
SELECT id, 'Samuel López' FROM courses WHERE name = 'Segundo A';
`)
	return err
}
