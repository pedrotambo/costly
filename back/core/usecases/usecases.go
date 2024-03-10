package usecases

import (
	"costly/core/ports/clock"
	"costly/core/ports/rpst"
	"errors"
)

type UseCases interface {
	IngredientUseCases
	RecipeUseCases
}

type useCases struct {
	IngredientUseCases
	RecipeUseCases
}

var ErrNotFound = errors.New("entity not found")
var ErrBadOpts = errors.New("bad create entity options")

func New(repository rpst.Repository, clock clock.Clock) UseCases {
	return &useCases{
		NewIngredientUseCases(repository, clock),
		NewRecipeUseCases(repository, clock),
	}
}
