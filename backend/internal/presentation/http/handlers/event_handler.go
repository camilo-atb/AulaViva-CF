package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/camilo-atb/aulaviva/backend/internal/application/dto"
	"github.com/camilo-atb/aulaviva/backend/internal/application/services"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/responses"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(service *services.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		events, err := h.service.List(r.Context())
		if err != nil {
			responses.Error(w, err)
			return
		}
		responses.JSON(w, http.StatusOK, events)
	case http.MethodPost:
		var input dto.CreateEventInput
		if err := decodeJSON(r, &input); err != nil {
			responses.Error(w, err)
			return
		}
		event, err := h.service.Create(r.Context(), input)
		if err != nil {
			responses.Error(w, err)
			return
		}
		responses.JSON(w, http.StatusCreated, event)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r.URL.Path, "/api/events/")
	if err != nil {
		responses.Error(w, err)
		return
	}
	event, err := h.service.Get(r.Context(), id)
	if err != nil {
		responses.Error(w, err)
		return
	}
	responses.JSON(w, http.StatusOK, event)
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return domainerrors.ErrInvalidInput
	}
	return nil
}

func pathID(path, prefix string) (int64, error) {
	value := strings.TrimPrefix(path, prefix)
	if value == path || strings.Contains(value, "/") {
		return 0, domainerrors.ErrInvalidInput
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, domainerrors.ErrInvalidInput
	}
	return id, nil
}
