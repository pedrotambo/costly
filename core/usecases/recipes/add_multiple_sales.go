package recipes

import (
	"context"
	"costly/core/errs"
	"costly/core/model"
	repo "costly/core/ports/repository"
)

type RecipeMultipleSalesAdder interface {
	AddMultipleSales(ctx context.Context, opts []RecipeSalesOptions) ([]*model.RecipeSales, error)
}

type RecipeSalesOptions struct {
	RecipeID  int64
	SoldUnits int
}

func (cr *recipeUseCases) AddMultipleSales(ctx context.Context, opts []RecipeSalesOptions) ([]*model.RecipeSales, error) {
	for _, sale := range opts {
		if sale.SoldUnits <= 0 {
			return []*model.RecipeSales{}, errs.ErrBadStockUnits
		}
	}

	allRecipeSales := []*model.RecipeSales{}
	cr.repository.Atomic(ctx, func(repo repo.Repository) error {
		for _, sale := range opts {
			recipeSales := model.NewRecipeSales(sale.RecipeID, sale.SoldUnits, cr.clock.Now())
			if err := repo.RecipeSales().Add(ctx, recipeSales); err != nil {
				return err
			}
			recipeIngredients, err := repo.Recipes().FindIngredients(ctx, recipeSales.RecipeID)
			if err != nil {
				return err
			}
			for _, recipeIngredient := range recipeIngredients {
				if err := repo.Ingredients().DecreaseStock(ctx, recipeIngredient.ID, recipeSales.Units*recipeIngredient.Units, recipeSales.CreatedAt); err != nil {
					return err
				}
			}
			allRecipeSales = append(allRecipeSales, recipeSales)
		}
		return nil
	})
	return allRecipeSales, nil
}
