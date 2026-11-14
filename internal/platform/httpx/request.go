package httpx

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/charmingruby/winston/internal/platform/validator"
)

var (
	ErrInvalidPayload = errors.New("invalid payload")
	ErrMissingParam   = errors.New("missing param")
)

func ParseRequest[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var obj T

	if err := json.UnmarshalRead(r.Body, &obj); err != nil {
		writeResponse(w, http.StatusBadRequest, ErrorResponse{
			Message: ErrInvalidPayload.Error(),
		})

		return nil, err
	}

	if err := validator.Validate(obj); err != nil {
		HandleErr(w, r, err)

		return nil, err
	}

	return &obj, nil
}

func GetPathParam(r *http.Request, key string) (string, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingParam, key)
	}

	return param, nil
}

func GetQueryParam(r *http.Request, key string) (string, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingParam, key)
	}

	return param, nil
}
