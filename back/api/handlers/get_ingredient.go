package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"net/http"
	"strconv"
)

func GetIngredientHandler(ingredientRepository rpst.IngredientRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDstr := r.PathValue("ingredientID")
		ingredientID, err := strconv.ParseInt(ingredientIDstr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
		}

		ingredient, err := ingredientRepository.GetIngredient(r.Context(), ingredientID)
		if err == rpst.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(r.Context(), err, "error getting ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		RespondJSON(w, 200, ingredient)
	}
}
