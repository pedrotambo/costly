package handlers

import (
	"costly/core/components/recipes"
	"costly/core/errs"
	"costly/core/model"
	"costly/core/ports/logger"
	"net/http"
	"strconv"
)

type RecipeResponse struct {
	model.Recipe
	Cost float64 `json:"cost"`
}

func GetRecipeHandler(recipeGetter recipes.RecipeFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipeIDstr := r.PathValue("recipeID")
		recipeID, err := strconv.ParseInt(recipeIDstr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
		}
		recipe, err := recipeGetter.Find(r.Context(), recipeID)
		if err == errs.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error(r.Context(), err, "error getting recipe")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		RespondJSON(w, 200, RecipeResponse{
			Recipe: recipe,
			Cost:   recipe.Cost(),
		})
	}
}
