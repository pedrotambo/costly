package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"net/http"
	"strconv"
)

func parseIngredientStockOptions(r *http.Request) (usecases.IngredientStockOptions, error) {
	opts := usecases.IngredientStockOptions{}
	if err := UnmarshallJSONBody(r, &opts); err != nil {
		return usecases.IngredientStockOptions{}, ErrBadJson
	}

	if opts.Units <= 0 {
		return usecases.IngredientStockOptions{}, ErrBadStockUnits
	}

	if opts.Price <= 0 {
		return usecases.IngredientStockOptions{}, ErrBadPrice
	}

	return opts, nil
}

func AddIngredientStockHandler(ingredientStockAdder usecases.IngredientStockAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDstr := r.PathValue("ingredientID")
		ingredientID, err := strconv.ParseInt(ingredientIDstr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
			return
		}

		ingredientStockOptions, err := parseIngredientStockOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
			return
		}

		ingredientStock, err := ingredientStockAdder.AddIngredientStock(r.Context(), int64(ingredientID), ingredientStockOptions)
		if err == rpst.ErrNotFound {
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
