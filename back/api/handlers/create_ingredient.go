package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"net/http"
)

func parseIngredientOptions(r *http.Request) (rpst.CreateIngredientOptions, error) {
	createIngredientOpts := rpst.CreateIngredientOptions{}
	if err := UnmarshallJSONBody(r, &createIngredientOpts); err != nil {
		return rpst.CreateIngredientOptions{}, ErrBadJson
	}

	if createIngredientOpts.Name == "" {
		return rpst.CreateIngredientOptions{}, ErrBadName
	}

	if createIngredientOpts.Unit != "gr" {
		return rpst.CreateIngredientOptions{}, ErrBadUnit
	}

	if createIngredientOpts.Price <= 0 {
		return rpst.CreateIngredientOptions{}, ErrBadPrice
	}

	return createIngredientOpts, nil
}

func CreateIngredientHandler(ingredientRepository rpst.IngredientRepository) http.HandlerFunc {
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
