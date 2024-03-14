package recipes

import (
	"context"
	"costly/core/errs"
	"costly/core/model"
	repo "costly/core/ports/repository"
)

type RecipeSalesAdder interface {
	AddSales(ctx context.Context, recipeID int64, soldUnits int) (*model.RecipeSales, error)
}

func (cr *recipeUseCases) AddSales(ctx context.Context, recipeID int64, soldUnits int) (*model.RecipeSales, error) {
	if soldUnits <= 0 {
		return &model.RecipeSales{}, errs.ErrBadStockUnits
	}
	recipeSales := model.NewRecipeSales(recipeID, soldUnits, cr.clock.Now())
	cr.repository.Atomic(ctx, func(repo repo.Repository) error {
		if err := repo.RecipeSales().Add(ctx, recipeSales); err != nil {
			return err
		}
		recipeIngredients, err := repo.Recipes().FindIngredients(ctx, recipeID)
		if err != nil {
			return err
		}
		for _, recipeIngredient := range recipeIngredients {
			if err := repo.Ingredients().DecreaseStock(ctx, recipeIngredient.ID, recipeSales.Units*recipeIngredient.Units, recipeSales.CreatedAt); err != nil {
				return err
			}
		}
		return nil
	})
	return recipeSales, nil
}
