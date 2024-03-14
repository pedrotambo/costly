package ingredients

import (
	"costly/core/ports/clock"
	"costly/core/ports/database"
	repo "costly/core/ports/repository"
)

type IngredientUseCases interface {
	IngredientCreator
	IngredientEditor
	IngredientStockAdder
	IngredientFinder
	IngredientsFinder
}

type ingredientUseCases struct {
	clock      clock.Clock
	repository repo.Repository
}

func New(database database.Database, clock clock.Clock) IngredientUseCases {
	return &ingredientUseCases{
		clock:      clock,
		repository: repo.New(database),
	}
}
