package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/camilo-atb/aulaviva/backend/internal/application/dto"
	"github.com/camilo-atb/aulaviva/backend/internal/application/services"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/responses"
)

type AuthorizedPersonHandler struct {
	service *services.AuthorizedPersonService
}

func NewAuthorizedPersonHandler(service *services.AuthorizedPersonService) *AuthorizedPersonHandler {
	return &AuthorizedPersonHandler{service: service}
}

func (h *AuthorizedPersonHandler) Create(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[0] != "api" || parts[1] != "students" || parts[3] != "authorized-people" {
		responses.Error(w, domainerrors.ErrInvalidInput)
		return
	}
	studentID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || studentID <= 0 {
		responses.Error(w, domainerrors.ErrInvalidInput)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var input dto.AuthorizedPersonInput
	if err := decodeJSON(r, &input); err != nil {
		responses.Error(w, err)
		return
	}
	person, err := h.service.Create(r.Context(), studentID, input)
	if err != nil {
		responses.Error(w, err)
		return
	}
	responses.JSON(w, http.StatusCreated, person)
}

func (h *AuthorizedPersonHandler) UpdateOrDelete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r.URL.Path, "/api/authorized-people/")
	if err != nil {
		responses.Error(w, domainerrors.ErrInvalidInput)
		return
	}
	switch r.Method {
	case http.MethodPut:
		var input dto.AuthorizedPersonInput
		if err := decodeJSON(r, &input); err != nil {
			responses.Error(w, err)
			return
		}
		person, err := h.service.Update(r.Context(), id, input)
		if err != nil {
			responses.Error(w, err)
			return
		}
		responses.JSON(w, http.StatusOK, person)
	case http.MethodDelete:
		if err := h.service.Delete(r.Context(), id); err != nil {
			responses.Error(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
