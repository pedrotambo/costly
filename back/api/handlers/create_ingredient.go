package handlers

import (
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/errs"
	"errors"
	"net/http"
)

func CreateIngredientHandler(ingredientCreator ingredients.IngredientCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		createIngredientOpts := ingredients.CreateIngredientOptions{}
		if err := UnmarshallJSONBody(r, &createIngredientOpts); err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadJson)
			return
		}
		ingredient, err := ingredientCreator.Create(r.Context(), createIngredientOpts)
		if errors.Is(err, errs.ErrBadOpts) {
			RespondJSON(w, http.StatusBadRequest, NewInvalidInputResponseError(err.Error()))
			return
		} else if err != nil {
			logger.Error(r.Context(), err, "error getting ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		RespondJSON(w, http.StatusCreated, ingredient)
	}
}
