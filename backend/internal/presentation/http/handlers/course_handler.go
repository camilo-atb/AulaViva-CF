package handlers

import (
	"net/http"

	"github.com/camilo-atb/aulaviva/backend/internal/application/services"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/responses"
)

type CourseHandler struct {
	service *services.CourseService
}

func NewCourseHandler(service *services.CourseService) *CourseHandler {
	return &CourseHandler{service: service}
}

func (h *CourseHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	courses, err := h.service.List(r.Context())
	if err != nil {
		responses.Error(w, err)
		return
	}
	responses.JSON(w, http.StatusOK, courses)
}
