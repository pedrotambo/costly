package recipes

import (
	"context"
	"costly/core/model"
	repo "costly/core/ports/repository"
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
	if err := cr.repository.Atomic(ctx, func(repo repo.Repository) error {
		if err := repo.Recipes().Add(ctx, newRecipe); err != nil {
			return fmt.Errorf("failed to create recipe: %s", err)
		}
		return nil
	}); err != nil {
		return &model.Recipe{}, err
	}

	return newRecipe, nil
}
