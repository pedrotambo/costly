package handlers

import (
	"costly/core/ports/repository"
	"net/http"

	"github.com/rs/zerolog"
)

func GetIngredientsHandler(ingredientRepository repository.IngredientRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredients, err := ingredientRepository.GetIngredients(r.Context())
		if err == repository.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			zerolog.Ctx(r.Context()).Error().Err(err).Msg("error getting ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, 200, ingredients)
	}
}
