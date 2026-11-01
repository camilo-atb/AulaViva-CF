package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/camilo-atb/aulaviva/backend/internal/application/services"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/responses"
)

type StudentHandler struct {
	service *services.StudentService
}

func NewStudentHandler(service *services.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

func (h *StudentHandler) ListOrGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/students" {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		students, err := h.service.List(r.Context(), r.URL.Query().Get("search"))
		if err != nil {
			responses.Error(w, err)
			return
		}
		responses.JSON(w, http.StatusOK, students)
		return
	}
	id, err := pathID(r.URL.Path, "/api/students/")
	if err != nil {
		responses.Error(w, err)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	student, err := h.service.Get(r.Context(), id)
	if err != nil {
		responses.Error(w, err)
		return
	}
	responses.JSON(w, http.StatusOK, student)
}

func (h *StudentHandler) ListPeople(w http.ResponseWriter, r *http.Request) {
	const prefix = "/api/students/"
	value := strings.TrimPrefix(r.URL.Path, prefix)
	parts := strings.Split(strings.Trim(value, "/"), "/")
	if len(parts) != 2 || parts[1] != "authorized-people" {
		responses.Error(w, domainerrors.ErrInvalidInput)
		return
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		responses.Error(w, domainerrors.ErrInvalidInput)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	people, err := h.service.ListPeople(r.Context(), id)
	if err != nil {
		responses.Error(w, err)
		return
	}
	responses.JSON(w, http.StatusOK, people)
}
