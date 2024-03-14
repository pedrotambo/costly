package recipes

import (
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	repo "costly/core/ports/repository"
	"costly/core/usecases/ingredients"
)

type RecipeUseCases interface {
	RecipeCreator
	RecipeSalesAdder
	RecipeFinder
	RecipesFinder
}

type recipeUseCases struct {
	clock       clock.Clock
	ingredients ingredients.IngredientUseCases
	repository  repo.Repository
}

func New(database database.Database, clock clock.Clock, logger logger.Logger, ingredients ingredients.IngredientUseCases) RecipeUseCases {
	return &recipeUseCases{
		clock:       clock,
		ingredients: ingredients,
		repository:  repo.New(database),
	}
}
