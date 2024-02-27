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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		editReq := repository.CreateIngredientOptions{}
		if err := UnmarshallJSONBody(r, &editReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := validateIngredientOptions(editReq); err != nil {
			vErr, ok := err.(ValidationError)
			if ok {
				RespondJSON(w, http.StatusBadRequest, ValidationErrorResponse{
					Errors: []ValidationError{vErr},
				})
				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = ingredientRepository.EditIngredient(r.Context(), int64(ingredientID), editReq.Name, editReq.Price, editReq.Unit)
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
