package handlers

import (
	"costly/core/components/logger"
	"costly/core/components/recipes"
	"costly/core/errs"
	"errors"
	"net/http"
)

func CreateRecipeHandler(recipeCreator recipes.RecipeCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createRecipeOptions := recipes.CreateRecipeOptions{}
		if err := UnmarshallJSONBody(r, &createRecipeOptions); err != nil {
			RespondJSON(w, http.StatusBadRequest, ErrBadJson)
			return
		}
		recipe, err := recipeCreator.Create(r.Context(), createRecipeOptions)
		if errors.Is(err, errs.ErrBadOpts) {
			RespondJSON(w, http.StatusBadRequest, NewInvalidInputResponseError(err.Error()))
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
