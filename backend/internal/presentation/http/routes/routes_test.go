package routes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/camilo-atb/aulaviva/backend/internal/application/services"
	sqliteinfra "github.com/camilo-atb/aulaviva/backend/internal/infrastructure/persistence/sqlite"
	sqliterepositories "github.com/camilo-atb/aulaviva/backend/internal/infrastructure/persistence/sqlite/repositories"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/handlers"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/routes"
)

func testServer(t *testing.T) (http.Handler, func()) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	ctx := context.Background()
	db, err := sqliteinfra.Open(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	if err := sqliteinfra.Seed(ctx, db); err != nil {
		t.Fatal(err)
	}
	courses := sqliterepositories.NewCourseRepository(db)
	events := sqliterepositories.NewEventRepository(db)
	authorizations := sqliterepositories.NewAuthorizationRepository(db)
	students := sqliterepositories.NewStudentRepository(db)
	people := sqliterepositories.NewAuthorizedPersonRepository(db)
	handler := routes.New(
		handlers.NewEventHandler(services.NewEventService(courses, events)),
		handlers.NewAuthorizationHandler(services.NewAuthorizationService(authorizations)),
		handlers.NewCourseHandler(services.NewCourseService(courses)),
		handlers.NewStudentHandler(services.NewStudentService(students, people)),
		handlers.NewAuthorizedPersonHandler(services.NewAuthorizedPersonService(students, people)),
	)
	return handler, func() { db.Close(); os.Remove(path) }
}

func TestEventCreationAndResponseFlow(t *testing.T) {
	handler, cleanup := testServer(t)
	defer cleanup()

	create := httptest.NewRequest(http.MethodPost, "/api/events", strings.NewReader(`{"name":"Salida al museo","description":"Visita educativa","eventDate":"2026-11-20","courseId":1}`))
	create.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, create)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body = %s", createResponse.Code, createResponse.Body.String())
	}
	var created struct {
		ID             int `json:"id"`
		Authorizations []struct {
			RequestToken string `json:"requestToken"`
		} `json:"authorizations"`
	}
	decode(t, createResponse, &created)
	if len(created.Authorizations) != 2 {
		t.Fatalf("authorizations = %d, want 2", len(created.Authorizations))
	}

	detail := httptest.NewRequest(http.MethodGet, "/api/events/1", nil)
	detailResponse := httptest.NewRecorder()
	handler.ServeHTTP(detailResponse, detail)
	if detailResponse.Code != http.StatusOK {
		t.Fatalf("detail status = %d", detailResponse.Code)
	}
	var event struct {
		Summary struct {
			Pending int `json:"pending"`
		} `json:"summary"`
	}
	decode(t, detailResponse, &event)
	if event.Summary.Pending != 2 {
		t.Fatalf("pending = %d, want 2", event.Summary.Pending)
	}

	token := created.Authorizations[0].RequestToken
	response := httptest.NewRequest(http.MethodPost, "/api/requests/"+token+"/response", strings.NewReader(`{"status":"authorized"}`))
	response.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	handler.ServeHTTP(responseRecorder, response)
	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("response status = %d, body = %s", responseRecorder.Code, responseRecorder.Body.String())
	}

	duplicate := httptest.NewRequest(http.MethodPost, "/api/requests/"+token+"/response", strings.NewReader(`{"status":"not_authorized"}`))
	duplicate.Header.Set("Content-Type", "application/json")
	duplicateRecorder := httptest.NewRecorder()
	handler.ServeHTTP(duplicateRecorder, duplicate)
	if duplicateRecorder.Code != http.StatusConflict {
		t.Fatalf("duplicate status = %d, want %d", duplicateRecorder.Code, http.StatusConflict)
	}
}

func TestHTTPValidationAndNotFound(t *testing.T) {
	handler, cleanup := testServer(t)
	defer cleanup()

	invalid := httptest.NewRequest(http.MethodPost, "/api/events", strings.NewReader(`{"name":"","eventDate":"bad","courseId":1}`))
	invalid.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, invalid)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("invalid status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}

	notFound := httptest.NewRequest(http.MethodGet, "/api/requests/missing", nil)
	notFoundRecorder := httptest.NewRecorder()
	handler.ServeHTTP(notFoundRecorder, notFound)
	if notFoundRecorder.Code != http.StatusNotFound {
		t.Fatalf("not found status = %d, want %d", notFoundRecorder.Code, http.StatusNotFound)
	}
}

func TestAuthorizedPeopleCRUDAndSearch(t *testing.T) {
	handler, cleanup := testServer(t)
	defer cleanup()

	list := httptest.NewRecorder()
	handler.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/api/students?search=Sofía", nil))
	if list.Code != http.StatusOK {
		t.Fatalf("student search status = %d", list.Code)
	}
	var students []struct {
		ID               int   `json:"id"`
		AuthorizedPeople []any `json:"authorizedPeople"`
	}
	decode(t, list, &students)
	if len(students) != 1 || len(students[0].AuthorizedPeople) != 0 {
		t.Fatalf("student search result = %+v", students)
	}

	create := httptest.NewRequest(http.MethodPost, "/api/students/1/authorized-people",
		strings.NewReader(`{"fullName":"Laura Gómez","documentType":"CC","documentNumber":"123456789","relationship":"Tía","phone":"+57 300 555 0101"}`))
	create.Header.Set("Content-Type", "application/json")
	createdResponse := httptest.NewRecorder()
	handler.ServeHTTP(createdResponse, create)
	if createdResponse.Code != http.StatusCreated {
		t.Fatalf("person create status = %d, body = %s", createdResponse.Code, createdResponse.Body.String())
	}
	var person struct {
		ID int `json:"id"`
	}
	decode(t, createdResponse, &person)

	duplicate := httptest.NewRecorder()
	handler.ServeHTTP(duplicate, httptest.NewRequest(http.MethodPost, "/api/students/1/authorized-people",
		strings.NewReader(`{"fullName":"Laura Gómez","documentType":"CC","documentNumber":"123456789","relationship":"Tía","phone":"3005550101"}`)))
	if duplicate.Code != http.StatusConflict {
		t.Fatalf("duplicate status = %d, want %d", duplicate.Code, http.StatusConflict)
	}

	update := httptest.NewRequest(http.MethodPut, "/api/authorized-people/"+strconv.Itoa(person.ID),
		strings.NewReader(`{"fullName":"Laura Gómez Pérez","documentType":"CC","documentNumber":"123456789","relationship":"Madre","phone":"300 555 0102"}`))
	update.Header.Set("Content-Type", "application/json")
	updatedResponse := httptest.NewRecorder()
	handler.ServeHTTP(updatedResponse, update)
	if updatedResponse.Code != http.StatusOK {
		t.Fatalf("person update status = %d", updatedResponse.Code)
	}

	deleteResponse := httptest.NewRecorder()
	handler.ServeHTTP(deleteResponse, httptest.NewRequest(http.MethodDelete, "/api/authorized-people/"+strconv.Itoa(person.ID), nil))
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("person delete status = %d", deleteResponse.Code)
	}

	detail := httptest.NewRecorder()
	handler.ServeHTTP(detail, httptest.NewRequest(http.MethodGet, "/api/students/1", nil))
	if detail.Code != http.StatusOK {
		t.Fatalf("student detail status = %d", detail.Code)
	}
	var student struct {
		AuthorizedPeople []any `json:"authorizedPeople"`
	}
	decode(t, detail, &student)
	if len(student.AuthorizedPeople) != 0 {
		t.Fatalf("authorized people after delete = %d", len(student.AuthorizedPeople))
	}
}

func decode(t *testing.T, recorder *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(recorder.Body).Decode(target); err != nil {
		t.Fatal(err)
	}
}
