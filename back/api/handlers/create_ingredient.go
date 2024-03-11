package handlers

import (
	"costly/core/ports/logger"
	"costly/core/usecases"
	"net/http"
)

func parseIngredientOptions(r *http.Request) (usecases.CreateIngredientOptions, error) {
	createIngredientOpts := usecases.CreateIngredientOptions{}
	if err := UnmarshallJSONBody(r, &createIngredientOpts); err != nil {
		return usecases.CreateIngredientOptions{}, ErrBadJson
	}

	if createIngredientOpts.Name == "" {
		return usecases.CreateIngredientOptions{}, ErrBadName
	}

	if createIngredientOpts.Unit != "gr" {
		return usecases.CreateIngredientOptions{}, ErrBadUnit
	}

	if createIngredientOpts.Price <= 0 {
		return usecases.CreateIngredientOptions{}, ErrBadPrice
	}

	return createIngredientOpts, nil
}

func CreateIngredientHandler(ingredientCreator usecases.IngredientCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createIngredientOpts, err := parseIngredientOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
			return
		}

		ingredient, err := ingredientCreator.CreateIngredient(r.Context(), createIngredientOpts)
		if err != nil {
			logger.Error(r.Context(), err, "error getting ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, http.StatusCreated, ingredient)
	}
}
