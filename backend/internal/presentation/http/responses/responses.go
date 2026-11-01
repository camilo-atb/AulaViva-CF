package responses

import (
	"encoding/json"
	"errors"
	"net/http"

	domainerrors "github.com/camilo-atb/aulaviva/backend/internal/domain/errors"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func Error(w http.ResponseWriter, err error) {
	status, code, message := http.StatusInternalServerError, "internal_error", "Ocurrió un error inesperado."
	switch {
	case errors.Is(err, domainerrors.ErrInvalidInput), errors.Is(err, domainerrors.ErrInvalidStatus):
		status, code, message = http.StatusBadRequest, "invalid_input", "La información enviada no es válida."
	case errors.Is(err, domainerrors.ErrNotFound):
		status, code, message = http.StatusNotFound, "not_found", "El recurso solicitado no existe."
	case errors.Is(err, domainerrors.ErrAlreadyAnswered), errors.Is(err, domainerrors.ErrConflict):
		status, code, message = http.StatusConflict, "authorization_already_answered", "La solicitud ya tiene una respuesta."
	}
	JSON(w, status, ErrorResponse{Error: ErrorBody{Code: code, Message: message}})
}
