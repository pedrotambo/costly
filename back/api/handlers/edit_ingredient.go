package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"net/http"
	"strconv"
)

func EditIngredientHandler(ingredientRepository repository.IngredientRepository) http.HandlerFunc {
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

		_, err = ingredientRepository.EditIngredient(r.Context(), int64(ingredientID), editIngredientOpts.Name, editIngredientOpts.Price, editIngredientOpts.Unit)
		if err == repository.ErrNotFound {
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
