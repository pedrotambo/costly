package recipes

import (
	"context"
	"costly/core/model"
	"fmt"
)

type RecipeCreator interface {
	Create(ctx context.Context, recipeOpts CreateRecipeOptions) (*model.Recipe, error)
}

type CreateRecipeOptions struct {
	Name        string
	Ingredients []model.RecipeIngredient
}

func (cr *recipeUseCases) Create(ctx context.Context, recipeOpts CreateRecipeOptions) (*model.Recipe, error) {
	newRecipe, err := model.NewRecipe(recipeOpts.Name, recipeOpts.Ingredients, cr.clock.Now())
	if err != nil {
		return &model.Recipe{}, err
	}

	if err := cr.repository.Recipes().Add(ctx, newRecipe); err != nil {
		return &model.Recipe{}, fmt.Errorf("failed to create recipe: %s", err)
	}

	return newRecipe, nil
}
