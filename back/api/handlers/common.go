package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ValidationError struct {
	OffendingField string `json:"offending_field"`
	Suggestion     string `json:"suggestion"`
}

func NewValidationError(offendingField, suggestion string) ValidationError {
	return ValidationError{offendingField, suggestion}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid field %s: %s", e.OffendingField, e.Suggestion)
}

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

func RespondJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonData, err := json.Marshal(body)
	if err != nil {
		// zerolog.Ctx(r.Context()).Error().Err(err).Msg("error getting recipe")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
	// json.NewEncoder(w).Encode(body)
}

// UnmarshallJSONBody reads the request body and unmarshalls it into the given interface.
func UnmarshallJSONBody(r *http.Request, v interface{}) error {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bs, v); err != nil {
		return err
	}

	return nil
}
