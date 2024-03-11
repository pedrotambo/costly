package handlers

import (
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/errs"
	"errors"
	"net/http"
	"strconv"
)

func EditIngredientHandler(ingredientEditor ingredients.IngredientEditor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDstr := r.PathValue("ingredientID")
		ingredientID, err := strconv.ParseInt(ingredientIDstr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
			return
		}
		editIngredientOpts := ingredients.CreateIngredientOptions{}
		if err := UnmarshallJSONBody(r, &editIngredientOpts); err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadJson)
			return
		}
		err = ingredientEditor.Update(r.Context(), int64(ingredientID), editIngredientOpts)
		if errors.Is(err, errs.ErrBadOpts) {
			RespondJSON(w, http.StatusBadRequest, NewInvalidInputResponseError(err.Error()))
			return
		} else if err == errs.ErrNotFound {
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
