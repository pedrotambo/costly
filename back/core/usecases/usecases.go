package usecases

import (
	"costly/core/ports"
	"costly/core/usecases/ingredients"
	"costly/core/usecases/recipes"
)

type UseCases struct {
	Ingredients ingredients.IngredientUseCases
	Recipes     recipes.RecipeUseCases
}

func New(ports *ports.Ports) (*UseCases, error) {
	ingredientUseCases := ingredients.New(ports.Database, ports.Clock)
	return &UseCases{
		Ingredients: ingredientUseCases,
		Recipes:     recipes.New(ports.Database, ports.Clock, ports.Logger, ingredientUseCases),
	}, nil
}
