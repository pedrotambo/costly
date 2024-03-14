package ingredients

import (
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	repo "costly/core/ports/repository"
)

type IngredientComponent interface {
	IngredientCreator
	IngredientEditor
	IngredientStockAdder
	IngredientFinder
	IngredientsFinder
}

type ingredientComponent struct {
	clock      clock.Clock
	repository repo.Repository
}

func New(database database.Database, clock clock.Clock, logger logger.Logger) IngredientComponent {
	return &ingredientComponent{
		clock:      clock,
		repository: repo.New(database),
	}
}
