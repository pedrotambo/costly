package handlers

import (
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"net/http"
)

func GetIngredientsHandler(ingredientsGetter ingredients.IngredientsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredients, err := ingredientsGetter.GetIngredients(r.Context())
		if err != nil {
			logger.Error(r.Context(), err, "error getting ingredients")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, 200, ingredients)
	}
}
