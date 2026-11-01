package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/camilo-atb/aulaviva/backend/internal/application/services"
	sqliteinfra "github.com/camilo-atb/aulaviva/backend/internal/infrastructure/persistence/sqlite"
	sqliterepositories "github.com/camilo-atb/aulaviva/backend/internal/infrastructure/persistence/sqlite/repositories"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/handlers"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/middleware"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/routes"
)

func main() {
	ctx := context.Background()
	dbPath := getenv("AULAVIVA_DB_PATH", "aulaviva.db")
	port := getenv("PORT", getenv("AULAVIVA_PORT", "8080"))
	allowedOrigin := getenv("AULAVIVA_ALLOWED_ORIGIN", "http://localhost:4200")
	if directory := filepath.Dir(dbPath); directory != "." {
		if err := os.MkdirAll(directory, 0o750); err != nil {
			log.Fatal(err)
		}
	}
	db, err := sqliteinfra.Open(ctx, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := sqliteinfra.Seed(ctx, db); err != nil {
		log.Fatal(err)
	}

	courseRepository := sqliterepositories.NewCourseRepository(db)
	eventRepository := sqliterepositories.NewEventRepository(db)
	authorizationRepository := sqliterepositories.NewAuthorizationRepository(db)
	studentRepository := sqliterepositories.NewStudentRepository(db)
	authorizedPersonRepository := sqliterepositories.NewAuthorizedPersonRepository(db)
	eventService := services.NewEventService(courseRepository, eventRepository)
	courseService := services.NewCourseService(courseRepository)
	authorizationService := services.NewAuthorizationService(authorizationRepository)
	studentService := services.NewStudentService(studentRepository, authorizedPersonRepository)
	authorizedPersonService := services.NewAuthorizedPersonService(studentRepository, authorizedPersonRepository)
	handler := routes.New(
		handlers.NewEventHandler(eventService),
		handlers.NewAuthorizationHandler(authorizationService),
		handlers.NewCourseHandler(courseService),
		handlers.NewStudentHandler(studentService),
		handlers.NewAuthorizedPersonHandler(authorizedPersonService),
	)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           middleware.CORS(allowedOrigin, handler),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("AulaViva API listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
