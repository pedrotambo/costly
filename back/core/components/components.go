package comps

import (
	"costly/core/components/clock"
	"costly/core/components/database"
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/components/recipes"
	"fmt"
)

type Components struct {
	Ingredients ingredients.IngredientComponent
	Recipes     recipes.RecipeComponent
	Clock       clock.Clock
	Logger      logger.Logger
	Database    database.Database
}

type Config struct {
	LogLevel string
	Database struct {
		ConnectionString string
	}
}

func InitComponents(config *Config) (*Components, error) {
	logger, err := logger.New(config.LogLevel)
	if err != nil {
		fmt.Printf("Could not create logger. Err: %s\n", err)
		return &Components{}, err
	}
	database, err := database.New(config.Database.ConnectionString, logger)
	if err != nil {
		logger.Error(err, "could not initialize database")
		return &Components{}, err
	}
	clock := clock.New()
	ingredientComponent := ingredients.New(database, clock, logger)
	recipeComponent := recipes.New(database, clock, logger, ingredientComponent)
	return &Components{
		Ingredients: ingredientComponent,
		Recipes:     recipeComponent,
		Logger:      logger,
		Database:    database,
		Clock:       clock,
	}, nil
}
