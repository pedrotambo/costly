package usecases

import (
	"costly/core/ports/clock"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"errors"
)

type Ports struct {
	Logger     logger.Logger
	Repository rpst.Repository
	Clock      clock.Clock
}

type UseCases interface {
	IngredientUseCases
	RecipeUseCases
	rpst.IngredientGetter
	rpst.IngredientsGetter
	rpst.IngredientStockUpdater
	rpst.RecipeGetter
	rpst.RecipesGetter
}

type useCases struct {
	IngredientUseCases
	RecipeUseCases
	rpst.Repository
}

var ErrNotFound = errors.New("entity not found")
var ErrBadOpts = errors.New("bad create entity options")

func New(ports *Ports) UseCases {
	return &useCases{
		NewIngredientUseCases(ports.Repository, ports.Clock),
		NewRecipeUseCases(ports.Repository, ports.Clock),
		ports.Repository,
	}
}
