package routes

import (
	"net/http"

	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/handlers"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/middleware"
)

func New(eventHandler *handlers.EventHandler, authorizationHandler *handlers.AuthorizationHandler, courseHandler *handlers.CourseHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/events" {
			eventHandler.ListOrCreate(w, r)
			return
		}
		if r.Method == http.MethodGet {
			eventHandler.Get(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/api/events/", eventHandler.Get)
	mux.HandleFunc("/api/requests/", authorizationHandler.GetOrRespond)
	mux.HandleFunc("/api/courses", courseHandler.List)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	return middleware.LocalCORS(mux)
}
