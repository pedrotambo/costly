package comps

import (
	"costly/core/components/clock"
	"costly/core/components/database"
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/components/recipes"
)

type Components struct {
	ingredients.IngredientComponent
	recipes.RecipeComponent
	clock.Clock
	logger.Logger
	database.Database
}
