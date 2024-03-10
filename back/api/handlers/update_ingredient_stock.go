package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"net/http"
	"strconv"
)

func parseNewStockOptions(r *http.Request) (rpst.NewStockOptions, error) {
	opts := rpst.NewStockOptions{}
	if err := UnmarshallJSONBody(r, &opts); err != nil {
		return rpst.NewStockOptions{}, ErrBadJson
	}

	if opts.NewUnits <= 0 {
		return rpst.NewStockOptions{}, ErrBadNewUnits
	}

	if opts.Price <= 0 {
		return rpst.NewStockOptions{}, ErrBadPrice
	}

	return opts, nil
}

func UpdateIngredientStockHandler(ingredientStockUpdater rpst.IngredientStockUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDstr := r.PathValue("ingredientID")
		ingredientID, err := strconv.ParseInt(ingredientIDstr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
			return
		}

		newStockOpts, err := parseNewStockOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
			return
		}

		_, err = ingredientStockUpdater.UpdateStock(r.Context(), int64(ingredientID), newStockOpts)
		if err == rpst.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(r.Context(), err, "error updating ingredient stock")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
