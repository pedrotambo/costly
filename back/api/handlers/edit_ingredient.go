package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"net/http"
	"strconv"
)

func EditIngredientHandler(ingredientEditor usecases.IngredientEditor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDstr := r.PathValue("ingredientID")
		ingredientID, err := strconv.ParseInt(ingredientIDstr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
			return
		}

		editIngredientOpts, err := parseIngredientOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
			return
		}

		err = ingredientEditor.EditIngredient(r.Context(), int64(ingredientID), editIngredientOpts)
		if err == rpst.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(r.Context(), err, "error getting ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
