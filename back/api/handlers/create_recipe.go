package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"net/http"
)

func parseRecipeOptions(r *http.Request) (repository.CreateRecipeOptions, error) {
	createRecipeOpts := repository.CreateRecipeOptions{}
	if err := UnmarshallJSONBody(r, &createRecipeOpts); err != nil {
		return repository.CreateRecipeOptions{}, ErrBadJson
	}
	if createRecipeOpts.Name == "" {
		return repository.CreateRecipeOptions{}, ErrBadName
	}

	if len(createRecipeOpts.Ingredients) == 0 {
		return repository.CreateRecipeOptions{}, ErrBadIngrs
	}

	return createRecipeOpts, nil
}

func CreateRecipeHandler(recipeRepository repository.RecipeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createRecipeOptions, err := parseRecipeOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
			return
		}

		recipe, err := recipeRepository.CreateRecipe(r.Context(), createRecipeOptions)
		if err == repository.ErrBadOpts {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			logger.Error(r.Context(), err, "error creating recipe")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, 201, RecipeResponse{
			Recipe: recipe,
			Cost:   recipe.Cost(),
		})
	}
}
