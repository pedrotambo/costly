package ingredients

import (
	"context"
	"costly/core/model"
	repo "costly/core/ports/repository"
)

type IngredientStockOptions struct {
	Units int
	Price float64
}

type IngredientStockAdder interface {
	AddStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error)
}

func (ic *ingredientUseCases) AddStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error) {
	ingredientStock, err := model.NewIngredientStock(ingredientID, ingredientStockOpts.Units, ingredientStockOpts.Price, ic.clock.Now())
	if err != nil {
		return &model.IngredientStock{}, err
	}
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
