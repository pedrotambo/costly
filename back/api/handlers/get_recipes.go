package handlers

import (
	costs "costly/core/logic"
	"costly/core/ports/repository"
	"net/http"

	"github.com/rs/zerolog"
)

func GetRecipesHandler(recipeRepository repository.RecipeRepository, costService costs.CostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipes, err := recipeRepository.GetRecipes(r.Context())
		if err == repository.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			zerolog.Ctx(r.Context()).Error().Err(err).Msg("error getting recipe")
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
