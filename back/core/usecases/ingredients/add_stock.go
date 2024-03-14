package ingredients

import (
	"context"
	"costly/core/errs"
	"costly/core/model"
	repo "costly/core/ports/repository"
)

type IngredientStockOptions struct {
	Units int
	Price float64
}

func (opts IngredientStockOptions) validate() error {
	if opts.Units <= 0 {
		return errs.ErrBadStockUnits
	}

	if opts.Price <= 0 {
		return errs.ErrBadPrice
	}
	return nil
}

type IngredientStockAdder interface {
	AddStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error)
}

func (ic *ingredientUseCases) AddStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error) {
	if err := ingredientStockOpts.validate(); err != nil {
		return &model.IngredientStock{}, err
	}
	ingredientStock := model.NewIngredientStock(ingredientID, ingredientStockOpts.Units, ingredientStockOpts.Price, ic.clock.Now())
	if err := (ic.repository.Atomic(ctx, func(repo repo.Repository) error {
		if err := repo.IngredientStocks().Add(ctx, ingredientStock); err != nil {
			return err
		}
		if err := repo.Ingredients().IncreaseStockAndUpdatePrice(ctx, ingredientID, ingredientStock.Units, ingredientStock.Price, ingredientStock.CreatedAt); err != nil {
			return err
		}
		return nil
	})); err != nil {
		return &model.IngredientStock{}, err
	}
	return ingredientStock, nil
}
