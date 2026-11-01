package repositories_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/camilo-atb/aulaviva/backend/internal/domain/entities"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
	sqliteinfra "github.com/camilo-atb/aulaviva/backend/internal/infrastructure/persistence/sqlite"
	repositories "github.com/camilo-atb/aulaviva/backend/internal/infrastructure/persistence/sqlite/repositories"
)

func TestPersistenceAndAtomicResponse(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "persistent.db")
	db, err := sqliteinfra.Open(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	if err := sqliteinfra.Seed(ctx, db); err != nil {
		t.Fatal(err)
	}
	courses := repositories.NewCourseRepository(db)
	course, err := courses.GetByID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	students, err := courses.ListStudents(ctx, course.ID)
	if err != nil {
		t.Fatal(err)
	}
	eventRepo := repositories.NewEventRepository(db)
	event, authorizations, err := eventRepo.CreateWithAuthorizations(ctx,
		entities.Event{CourseID: course.ID, Name: "Evento", EventDate: "2026-12-01"},
		[]entities.StudentAuthorization{
			{Student: students[0], Authorization: entities.Authorization{RequestToken: "token-a"}},
			{Student: students[1], Authorization: entities.Authorization{RequestToken: "token-b"}},
		})
	if err != nil {
		t.Fatal(err)
	}
	if len(authorizations) != 2 {
		t.Fatalf("authorizations = %d, want 2", len(authorizations))
	}
	authRepo := repositories.NewAuthorizationRepository(db)
	if err := authRepo.Respond(ctx, "token-a", entities.StatusAuthorized, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	if err := authRepo.Respond(ctx, "token-a", entities.StatusNotAuthorized, time.Now().UTC()); err != domainerrors.ErrAlreadyAnswered {
		t.Fatalf("duplicate error = %v, want %v", err, domainerrors.ErrAlreadyAnswered)
	}
	db.Close()

	db, err = sqliteinfra.Open(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	detail, _, err := repositories.NewEventRepository(db).GetDetails(ctx, event.ID)
	if err != nil {
		t.Fatal(err)
	}
	if detail.ID != event.ID {
		t.Fatalf("event id after reopen = %d, want %d", detail.ID, event.ID)
	}
}

func TestCreationKeepsUniqueEventStudentConstraint(t *testing.T) {
	ctx := context.Background()
	db, err := sqliteinfra.Open(ctx, filepath.Join(t.TempDir(), "unique.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := sqliteinfra.Seed(ctx, db); err != nil {
		t.Fatal(err)
	}
	courseRepo := repositories.NewCourseRepository(db)
	students, err := courseRepo.ListStudents(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	repo := repositories.NewEventRepository(db)
	event := entities.Event{CourseID: 1, Name: "Duplicado", EventDate: "2026-12-02"}
	items := []entities.StudentAuthorization{{Student: students[0], Authorization: entities.Authorization{RequestToken: "one"}}}
	created, _, err := repo.CreateWithAuthorizations(ctx, event, items)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx, `INSERT INTO authorizations (event_id, student_id, request_token, status) VALUES (?, ?, ?, 'pending')`, created.ID, students[0].ID, "two")
	if err == nil {
		t.Fatal("expected unique event-student constraint error")
	}
	_, _, err = repo.GetDetails(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
}
