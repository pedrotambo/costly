package handlers

import (
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"net/http"
)

func parseIngredientOptions(r *http.Request) (ingredients.CreateIngredientOptions, error) {
	createIngredientOpts := ingredients.CreateIngredientOptions{}
	if err := UnmarshallJSONBody(r, &createIngredientOpts); err != nil {
		return ingredients.CreateIngredientOptions{}, ErrBadJson
	}

	if createIngredientOpts.Name == "" {
		return ingredients.CreateIngredientOptions{}, ErrBadName
	}

	if createIngredientOpts.Unit != "gr" {
		return ingredients.CreateIngredientOptions{}, ErrBadUnit
	}

	if createIngredientOpts.Price <= 0 {
		return ingredients.CreateIngredientOptions{}, ErrBadPrice
	}

	return createIngredientOpts, nil
}

func CreateIngredientHandler(ingredientCreator ingredients.IngredientCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createIngredientOpts, err := parseIngredientOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
			return
		}

		ingredient, err := ingredientCreator.Create(r.Context(), createIngredientOpts)
		if err != nil {
			logger.Error(r.Context(), err, "error getting ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, http.StatusCreated, ingredient)
	}
}
