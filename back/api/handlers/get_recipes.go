package handlers

import (
	"costly/core/components/logger"
	"costly/core/components/recipes"
	"net/http"
)

func GetRecipesHandler(recipesGetter recipes.RecipesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipes, err := recipesGetter.GetRecipes(r.Context())
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
