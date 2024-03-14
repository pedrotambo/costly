package ingredients

import (
	"context"
	"costly/core/model"
)

type IngredientCreator interface {
	Create(ctx context.Context, ingredientOpts CreateIngredientOptions) (*model.Ingredient, error)
}

type CreateIngredientOptions struct {
	Name  string
	Price float64
	Unit  model.Unit
}

func (ic *ingredientComponent) Create(ctx context.Context, opts CreateIngredientOptions) (*model.Ingredient, error) {
	newIngredient, err := model.NewIngredient(opts.Name, opts.Unit, opts.Price, ic.clock.Now())
	if err != nil {
		return &model.Ingredient{}, err
	}
	if err := ic.repository.Ingredients().Add(ctx, newIngredient); err != nil {
		return nil, err
	}

	return newIngredient, nil
}
