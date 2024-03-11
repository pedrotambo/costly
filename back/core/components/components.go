package comps

import (
	"costly/core/components/ingredients"
	"costly/core/components/recipes"
	"costly/core/ports"
)

type Components struct {
	Ingredients ingredients.IngredientComponent
	Recipes     recipes.RecipeComponent
}

func New(ports *ports.Ports) (*Components, error) {
	ingredientComponent := ingredients.New(ports.Database, ports.Clock, ports.Logger)
	recipeComponent := recipes.New(ports.Database, ports.Clock, ports.Logger, ingredientComponent)
	return &Components{
		Ingredients: ingredientComponent,
		Recipes:     recipeComponent,
	}, nil
}
