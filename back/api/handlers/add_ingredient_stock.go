package handlers

import (
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/errs"
	"net/http"
	"strconv"
)

func parseIngredientStockOptions(r *http.Request) (ingredients.IngredientStockOptions, error) {
	opts := ingredients.IngredientStockOptions{}
	if err := UnmarshallJSONBody(r, &opts); err != nil {
		return ingredients.IngredientStockOptions{}, ErrBadJson
	}

	if opts.Units <= 0 {
		return ingredients.IngredientStockOptions{}, ErrBadStockUnits
	}

	if opts.Price <= 0 {
		return ingredients.IngredientStockOptions{}, ErrBadPrice
	}

	return opts, nil
}

func AddIngredientStockHandler(ingredientStockAdder ingredients.IngredientStockAdder) http.HandlerFunc {
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

		ingredientStock, err := ingredientStockAdder.AddStock(r.Context(), int64(ingredientID), ingredientStockOptions)
		if err == errs.ErrNotFound {
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
