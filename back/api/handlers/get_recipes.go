package handlers

import (
	"costly/core/components/recipes"
	"costly/core/ports/logger"
	"net/http"
)

func GetRecipesHandler(recipesGetter recipes.RecipesFinder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipes, err := recipesGetter.FindAll(r.Context())
		if err != nil {
			logger.Error(r.Context(), err, "error getting recipes")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		recipeResponses := []RecipeResponse{}
		for _, recipe := range recipes {
			recipeResponses = append(recipeResponses, RecipeResponse{
				RecipeView: recipe,
				Cost:       recipe.Cost(),
			})
		}
		RespondJSON(w, 200, recipeResponses)
	}
}
