package handlers

import (
	"costly/core/components/recipes"
	"costly/core/errs"
	"costly/core/ports/logger"
	"errors"
	"net/http"
	"strconv"
)

type recipeSalesOpts struct {
	SoldUnits int `json:"sold_units"`
}

func AddRecipeSalesHandler(recipeSalesAddres recipes.RecipeSalesAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipeIDStr := r.PathValue("recipeID")
		recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
			return
		}
		opts := recipeSalesOpts{}
		if err := UnmarshallJSONBody(r, &opts); err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadJson)
			return
		}
		ingredientStock, err := recipeSalesAddres.AddSales(r.Context(), recipeID, opts.SoldUnits)
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
