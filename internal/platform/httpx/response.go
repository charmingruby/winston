package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charmingruby/winstom/internal/platform/o11y/logging"
	"github.com/charmingruby/winstom/internal/platform/validator"
)

type ErrorResponse struct {
	Message string                 `json:"message"`
	Errors  []validator.FieldError `json:"errors,omitempty"`
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
	writeResponse(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
}

func HandleErr(w http.ResponseWriter, r *http.Request, err error) {
	if validationErr, ok := errors.AsType[*validator.ValidationError](err); ok {
		writeResponse(w, http.StatusBadRequest, ErrorResponse{
			Message: validationErr.Error(),
			Errors:  validationErr.Fields,
		})
		return
	}

	log := logging.LoggerFromContext(r.Context())

	log.Error("internal server error", "error", err)

	writeResponse(w, http.StatusInternalServerError, ErrorResponse{
		Message: "Internal Server Error",
	})
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
