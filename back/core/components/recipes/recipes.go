package recipes

import (
	"costly/core/components/ingredients"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	repo "costly/core/ports/repository"
)

type RecipeComponent interface {
	RecipeCreator
	RecipeSalesAdder
	RecipeFinder
	RecipesFinder
}

type recipeUseCases struct {
	clock       clock.Clock
	ingredients ingredients.IngredientComponent
	repository  repo.Repository
}

func New(database database.Database, clock clock.Clock, logger logger.Logger, ingredients ingredients.IngredientComponent) RecipeComponent {
	return &recipeUseCases{
		clock:       clock,
		ingredients: ingredients,
		repository:  repo.New(database),
	}
}
