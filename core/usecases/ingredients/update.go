package ingredients

import (
	"context"
	"costly/core/errs"
	"costly/core/model"
)

type IngredientEditor interface {
	Update(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) error
}

func (opts CreateIngredientOptions) validate() error {
	if opts.Name == "" {
		return errs.ErrBadName
	}
	if opts.Unit != "gr" {
		return errs.ErrBadUnit
	}
	if opts.Price <= 0 {
		return errs.ErrBadPrice
	}
	return nil
}

func (ic *ingredientUseCases) Update(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) error {
	if err := ingredientOpts.validate(); err != nil {
		return err
	}
	err := ic.repository.Ingredients().Update(ctx, ingredientID, func(ingredient *model.Ingredient) error {
		ingredient.Name = ingredientOpts.Name
		ingredient.Price = ingredientOpts.Price
		ingredient.Unit = ingredientOpts.Unit
		ingredient.LastModified = ic.clock.Now()
		return nil
	})
	return err
}
