package handlers

import (
	"costly/core/domain"
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"net/http"
)

type validationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

type IngredientRequest struct {
	Name  string      `json:"name"`
	Unit  domain.Unit `json:"unit"`
	Price float64     `json:"price"`
}

var ErrBadName = NewValidationError("name", "el name debe ser valido")
var ErrBadUnit = NewValidationError("unit", "la unidad es inválida")

func (req IngredientRequest) Validate() error {
	if req.Name == "" {
		return ErrBadName
	}

	if req.Unit != "gr" {
		return ErrBadUnit
	}

	return nil
}

func CreateIngredientHandler(repository repository.IngredientRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createReq := IngredientRequest{}
		if err := UnmarshallJSONBody(r, &createReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := createReq.Validate()

		if err != nil {
			RespondJSON(w, http.StatusBadRequest, validationErrorResponse{
				Errors: []ValidationError{err.(ValidationError)},
			})
			return
		}

		ingredient, err := repository.CreateIngredient(r.Context(), createReq.Name, createReq.Price, createReq.Unit)
		if err != nil {
			logger.Error(r.Context(), err, "error getting ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, 201, ingredient)
	}
}
