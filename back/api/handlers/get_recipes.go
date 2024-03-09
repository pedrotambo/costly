package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"net/http"
)

func GetRecipesHandler(recipeRepository rpst.RecipeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipes, err := recipeRepository.GetRecipes(r.Context())
		if err != nil {
			logger.Error(r.Context(), err, "error getting recipes")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		recipeResponses := []RecipeResponse{}
		for _, recipe := range recipes {
			recipeResponses = append(recipeResponses, RecipeResponse{
				Recipe: recipe,
				Cost:   recipe.Cost(),
			})
		}

		RespondJSON(w, 200, recipeResponses)
	}
}
