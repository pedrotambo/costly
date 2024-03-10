package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"net/http"
)

func parseRecipeOptions(r *http.Request) (usecases.CreateRecipeOptions, error) {
	createRecipeOpts := usecases.CreateRecipeOptions{}
	if err := UnmarshallJSONBody(r, &createRecipeOpts); err != nil {
		return usecases.CreateRecipeOptions{}, ErrBadJson
	}
	if createRecipeOpts.Name == "" {
		return usecases.CreateRecipeOptions{}, ErrBadName
	}

	if len(createRecipeOpts.Ingredients) == 0 {
		return usecases.CreateRecipeOptions{}, ErrBadIngrs
	}

	return createRecipeOpts, nil
}

func CreateRecipeHandler(recipeCreator usecases.RecipeCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createRecipeOptions, err := parseRecipeOptions(r)
		if err != nil {
			RespondJSON(w, http.StatusBadRequest, err)
			return
		}

		recipe, err := recipeCreator.CreateRecipe(r.Context(), createRecipeOptions)
		if err == rpst.ErrBadOpts {
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
