package handlers

import (
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"net/http"
)

type createRecipeRequest struct {
	Name        string                             `json:"name"`
	Ingredients []repository.RecipeIngredientInput `json:"ingredients"`
}

func (req createRecipeRequest) Validate() error {
	if req.Name == "" {
		return NewValidationError("name", "el name debe ser valido")
	}

	if len(req.Ingredients) == 0 {
		return NewValidationError("ingredients", "la receta tiene que tener algún ingrediente")
	}

	return nil
}

func CreateRecipeHandler(recipeRepository repository.RecipeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createReq := createRecipeRequest{}
		if err := UnmarshallJSONBody(r, &createReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO: refactor this out
		if err := createReq.Validate(); err != nil {
			vErr, ok := err.(ValidationError)
			if ok {
				RespondJSON(w, http.StatusBadRequest, ValidationErrorResponse{
					Errors: []ValidationError{vErr},
				})
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		recipe, err := recipeRepository.CreateRecipe(r.Context(), createReq.Name, createReq.Ingredients)
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
