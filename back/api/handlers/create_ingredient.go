package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"net/http"
)

type validationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

var ErrBadName = NewValidationError("name", "el name debe ser valido")
var ErrBadUnit = NewValidationError("unit", "la unidad es inválida")

func validateIngredientOptions(opts repository.CreateIngredientOptions) error {
	if opts.Name == "" {
		return ErrBadName
	}

	if opts.Unit != "gr" {
		return ErrBadUnit
	}

	return nil
}

func CreateIngredientHandler(ingredientRepository repository.IngredientRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createIngredientOpts := repository.CreateIngredientOptions{}
		if err := UnmarshallJSONBody(r, &createIngredientOpts); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := validateIngredientOptions(createIngredientOpts)

		if err != nil {
			RespondJSON(w, http.StatusBadRequest, validationErrorResponse{
				Errors: []ValidationError{err.(ValidationError)},
			})
			return
		}

		ingredient, err := ingredientRepository.CreateIngredient(r.Context(), createIngredientOpts)
		if err != nil {
			logger.Error(r.Context(), err, "error getting ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, http.StatusCreated, ingredient)
	}
}
