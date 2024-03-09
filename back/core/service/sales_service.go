package service

import (
	"context"
	"costly/core/ports/rpst"
)

type SoldRecipes struct {
	RecipeID  int64 `json:"recipe_id"`
	SoldUnits int   `json:"sold_units"`
}

type SalesService interface {
	RegisterSoldRecipes(ctx context.Context, soldRecipes SoldRecipes) error
}

type salesService struct {
	repository rpst.Repository
}

func NewSalesService(repository rpst.Repository) SalesService {
	return &salesService{repository: repository}
}

func (s *salesService) RegisterSoldRecipes(ctx context.Context, soldRecipes SoldRecipes) error {
	recipe, err := s.repository.GetRecipe(ctx, soldRecipes.RecipeID)

	if err != nil {
		return err
	}

	unitsUsedByID := map[int64]int{}
	for _, recipeIngredient := range recipe.Ingredients {
		unitsUsedByID[recipeIngredient.Ingredient.ID] = recipeIngredient.Units * soldRecipes.SoldUnits
	}

	// err = s.repository.ReduceStockUnits(ctx, unitsUsedByID)

	if err != nil {
		return err
	}
	return nil
}
