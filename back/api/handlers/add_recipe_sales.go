package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"net/http"
	"strconv"
)

func parseRecipeSalesOptions(r *http.Request) (recipeSalesOpts, error) {
	opts := recipeSalesOpts{}
	if err := UnmarshallJSONBody(r, &opts); err != nil {
		return recipeSalesOpts{}, ErrBadJson
	}

	if opts.SoldUnits <= 0 {
		return recipeSalesOpts{}, ErrBadStockUnits
	}

	return opts, nil
}

type recipeSalesOpts struct {
	SoldUnits int `json:"sold_units"`
}

func AddRecipeSalesHandler(recipeSalesAddres usecases.RecipeSalesAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipeIDStr := r.PathValue("recipeID")
		recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
			return
		}

		opts, err := parseRecipeSalesOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
			return
		}

		ingredientStock, err := recipeSalesAddres.AddRecipeSales(r.Context(), recipeID, opts.SoldUnits)
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
