package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"net/http"
)

func GetIngredientsHandler(ingredientRepository rpst.IngredientRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredients, err := ingredientRepository.GetIngredients(r.Context())
		if err != nil {
			logger.Error(r.Context(), err, "error getting ingredients")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, 200, ingredients)
	}
}
