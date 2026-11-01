package routes

import (
	"net/http"
	"strings"

	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/handlers"
)

func New(eventHandler *handlers.EventHandler, authorizationHandler *handlers.AuthorizationHandler, courseHandler *handlers.CourseHandler, studentHandler *handlers.StudentHandler, authorizedPersonHandler *handlers.AuthorizedPersonHandler) http.Handler {
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
	mux.HandleFunc("/api/students", studentHandler.ListOrGet)
	mux.HandleFunc("/api/students/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(strings.TrimSuffix(r.URL.Path, "/"), "/authorized-people") {
			if r.Method == http.MethodPost {
				authorizedPersonHandler.Create(w, r)
			} else {
				studentHandler.ListPeople(w, r)
			}
			return
		}
		studentHandler.ListOrGet(w, r)
	})
	mux.HandleFunc("/api/authorized-people/", authorizedPersonHandler.UpdateOrDelete)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}
