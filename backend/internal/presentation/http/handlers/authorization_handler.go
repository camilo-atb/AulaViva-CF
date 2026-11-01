package handlers

import (
	"net/http"
	"strings"

	"github.com/camilo-atb/aulaviva/backend/internal/application/dto"
	"github.com/camilo-atb/aulaviva/backend/internal/application/services"
	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
	"github.com/camilo-atb/aulaviva/backend/internal/presentation/http/responses"
)

type AuthorizationHandler struct {
	service *services.AuthorizationService
}

func NewAuthorizationHandler(service *services.AuthorizationService) *AuthorizationHandler {
	return &AuthorizationHandler{service: service}
}

func (h *AuthorizationHandler) GetOrRespond(w http.ResponseWriter, r *http.Request) {
	token, action, ok := requestPath(r.URL.Path)
	if !ok {
		responses.Error(w, domainerrors.ErrInvalidInput)
		return
	}
	if action == "response" {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var input dto.SubmitResponseInput
		if err := decodeJSON(r, &input); err != nil {
			responses.Error(w, err)
			return
		}
		result, err := h.service.Respond(r.Context(), token, input)
		if err != nil {
			responses.Error(w, err)
			return
		}
		responses.JSON(w, http.StatusOK, result)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	result, err := h.service.Get(r.Context(), token)
	if err != nil {
		responses.Error(w, err)
		return
	}
	responses.JSON(w, http.StatusOK, result)
}

func requestPath(path string) (string, string, bool) {
	const prefix = "/api/requests/"
	value := strings.TrimPrefix(path, prefix)
	if value == path {
		return "", "", false
	}
	parts := strings.Split(strings.Trim(value, "/"), "/")
	if len(parts) == 1 && parts[0] != "" {
		return parts[0], "", true
	}
	if len(parts) == 2 && parts[1] == "response" && parts[0] != "" {
		return parts[0], "response", true
	}
	return "", "", false
}
