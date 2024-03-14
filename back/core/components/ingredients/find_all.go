package ingredients

import (
	"context"
	"costly/core/model"
)

type IngredientsFinder interface {
	FindAll(ctx context.Context) ([]model.Ingredient, error)
}

func (ic *ingredientComponent) FindAll(ctx context.Context) ([]model.Ingredient, error) {
	return ic.repository.Ingredients().FindAll(ctx)
}
