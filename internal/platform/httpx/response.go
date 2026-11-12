package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charmingruby/stan/internal/platform/validator"
)

type ErrorResponse struct {
	Message string                 `json:"message"`
	Fields  []validator.FieldError `json:"fields,omitempty"`
}

func WriteOKResponse(w http.ResponseWriter, v any) {
	writeResponse(w, http.StatusOK, v)
}

func WriteCreatedResponse(w http.ResponseWriter, v any) {
	writeResponse(w, http.StatusCreated, v)
}

func WriteEmptyResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func WriteBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	response := ErrorResponse{
		Message: err.Error(),
	}
	if validationErr, ok := errors.AsType[*validator.ValidationError](err); ok {
		response.Fields = validationErr.Fields
	}

	writeResponse(w, http.StatusBadRequest, response)
}

func writeResponse(w http.ResponseWriter, status int, v any) {
	if v == nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"Internal Server Error"}`))
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(body)
}
