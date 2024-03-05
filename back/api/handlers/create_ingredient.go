package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"net/http"
)

func parseIngredientOptions(r *http.Request) (repository.CreateIngredientOptions, error) {
	createIngredientOpts := repository.CreateIngredientOptions{}
	if err := UnmarshallJSONBody(r, &createIngredientOpts); err != nil {
		return repository.CreateIngredientOptions{}, ErrBadJson
	}

	if createIngredientOpts.Name == "" {
		return repository.CreateIngredientOptions{}, ErrBadName
	}

	if createIngredientOpts.Unit != "gr" {
		return repository.CreateIngredientOptions{}, ErrBadUnit
	}

	if createIngredientOpts.Price <= 0 {
		return repository.CreateIngredientOptions{}, ErrBadPrice
	}

	return createIngredientOpts, nil
}

func CreateIngredientHandler(ingredientRepository repository.IngredientRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createIngredientOpts, err := parseIngredientOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
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
