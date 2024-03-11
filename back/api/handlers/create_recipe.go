package handlers

import (
	"costly/core/components/logger"
	"costly/core/components/recipes"
	"costly/core/errs"
	"net/http"
)

func parseRecipeOptions(r *http.Request) (recipes.CreateRecipeOptions, error) {
	createRecipeOpts := recipes.CreateRecipeOptions{}
	if err := UnmarshallJSONBody(r, &createRecipeOpts); err != nil {
		return recipes.CreateRecipeOptions{}, ErrBadJson
	}
	if createRecipeOpts.Name == "" {
		return recipes.CreateRecipeOptions{}, ErrBadName
	}

	if len(createRecipeOpts.Ingredients) == 0 {
		return recipes.CreateRecipeOptions{}, ErrBadIngrs
	}

	return createRecipeOpts, nil
}

func CreateRecipeHandler(recipeCreator recipes.RecipeCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createRecipeOptions, err := parseRecipeOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
			return
		}

		recipe, err := recipeCreator.CreateRecipe(r.Context(), createRecipeOptions)
		if err == errs.ErrBadOpts {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			logger.Error(r.Context(), err, "error creating recipe")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, 201, RecipeResponse{
			Recipe: *recipe,
			Cost:   recipe.Cost(),
		})
	}
}
