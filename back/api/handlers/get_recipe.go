package handlers

import (
	"costly/core/domain"
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RecipeResponse struct {
	domain.Recipe
	Cost float64 `json:"cost"`
}

func GetRecipeHandler(recipeRepository repository.RecipeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipeIDstr := chi.URLParam(r, "recipeID")
		recipeID, err := strconv.ParseInt(recipeIDstr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
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
