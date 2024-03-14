package ingredients

import (
	"context"
	"costly/core/model"
)

type IngredientFinder interface {
	Find(ctx context.Context, id int64) (model.Ingredient, error)
}

func (ic *ingredientComponent) Find(ctx context.Context, id int64) (model.Ingredient, error) {
	return ic.repository.Ingredients().Find(ctx, id)
}
