package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/camilo-atb/aulaviva/backend/internal/application/services"
	sqliteinfra "github.com/camilo-atb/aulaviva/backend/internal/infrastructure/persistence/sqlite"
	sqliterepositories "github.com/camilo-atb/aulaviva/backend/internal/infrastructure/persistence/sqlite/repositories"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/handlers"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/routes"
)

func main() {
	ctx := context.Background()
	dbPath := getenv("AULAVIVA_DB_PATH", "aulaviva.db")
	port := getenv("AULAVIVA_PORT", "8080")
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
	eventService := services.NewEventService(courseRepository, eventRepository)
	courseService := services.NewCourseService(courseRepository)
	authorizationService := services.NewAuthorizationService(authorizationRepository)
	handler := routes.New(handlers.NewEventHandler(eventService), handlers.NewAuthorizationHandler(authorizationService), handlers.NewCourseHandler(courseService))

	server := &http.Server{Addr: ":" + port, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	log.Printf("AulaViva API listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
