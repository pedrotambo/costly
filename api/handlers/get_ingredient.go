package handlers

import (
	"costly/core/errs"
	"costly/core/ports/logger"
	"costly/core/usecases/ingredients"
	"net/http"
	"strconv"
)

func GetIngredientHandler(ingredientGetter ingredients.IngredientFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDstr := r.PathValue("ingredientID")
		ingredientID, err := strconv.ParseInt(ingredientIDstr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
		}
		ingredient, err := ingredientGetter.Find(r.Context(), ingredientID)
		if err == errs.ErrNotFound {
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
