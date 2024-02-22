package handlers

import (
	"costly/core/ports/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func GetIngredientHandler(ingredientRepository repository.IngredientRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientIDstr := chi.URLParam(r, "ingredientID")
		ingredientID, err := strconv.ParseInt(ingredientIDstr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ingredient, err := ingredientRepository.GetIngredient(r.Context(), ingredientID)
		if err == repository.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			zerolog.Ctx(r.Context()).Error().Err(err).Msg("error getting ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, 200, ingredient)
	}
}
