package handlers

import (
	"costly/core/domain"
	"costly/core/ports/repository"
	"net/http"

	"github.com/rs/zerolog"
)

type createIngredientRequest struct {
	Name  string      `json:"name"`
	Unit  domain.Unit `json:"unit"`
	Price float64     `json:"price"`
}

func (req createIngredientRequest) Validate() error {
	if req.Name == "" {
		return NewValidationError("name", "el name debe ser valido")
	}

	if req.Unit != "gr" {
		return NewValidationError("unit", "la unidad es inválida")
	}

	return nil
}

func CreateIngredientHandler(repository repository.IngredientRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		createReq := createIngredientRequest{}
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

		ingredient, err := repository.CreateIngredient(r.Context(), createReq.Name, createReq.Price, createReq.Unit)
		if err != nil {
			zerolog.Ctx(r.Context()).Error().Err(err).Msg("error creating ingredient")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondJSON(w, 201, ingredient)
	}
}
