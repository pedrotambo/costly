package handlers

import (
	"costly/core/errs"
	"costly/core/ports/logger"
	"costly/core/usecases/ingredients"
	"errors"
	"net/http"
	"strconv"
)

func AddIngredientStockHandler(ingredientStockAdder ingredients.IngredientStockAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDstr := r.PathValue("ingredientID")
		ingredientID, err := strconv.ParseInt(ingredientIDstr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
			return
		}
		ingredientStockOptions := ingredients.IngredientStockOptions{}
		if err := UnmarshallJSONBody(r, &ingredientStockOptions); err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadJson)
			return
		}
		ingredientStock, err := ingredientStockAdder.AddStock(r.Context(), int64(ingredientID), ingredientStockOptions)
		if errors.Is(err, errs.ErrBadOpts) {
			RespondJSON(w, http.StatusBadRequest, NewInvalidInputResponseError(err.Error()))
			return
		} else if err == errs.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(r.Context(), err, "error adding ingredient stock")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, http.StatusCreated, ingredientStock)
	}
}
