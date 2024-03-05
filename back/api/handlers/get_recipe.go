package handlers

import (
	"costly/core/domain"
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"net/http"
	"strconv"
)

type RecipeResponse struct {
	domain.Recipe
	Cost float64 `json:"cost"`
}

func GetRecipeHandler(recipeRepository repository.RecipeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipeIDstr := r.PathValue("recipeID")
		recipeID, err := strconv.ParseInt(recipeIDstr, 10, 64)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadID)
		}

		recipe, err := recipeRepository.GetRecipe(r.Context(), recipeID)
		if err == repository.ErrNotFound {
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
