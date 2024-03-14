package recipes

import (
	"context"
	"costly/core/model"
	repo "costly/core/ports/repository"
)

type RecipeFinder interface {
	Find(ctx context.Context, id int64) (model.RecipeView, error)
}

func (cr *recipeUseCases) Find(ctx context.Context, id int64) (model.RecipeView, error) {
	var recipe model.Recipe
	var recipeIngredientsView []model.RecipeIngredientView
	if err := cr.repository.Atomic(ctx, func(repo repo.Repository) error {
		recipeFound, err := repo.Recipes().Find(ctx, id)
		if err != nil {
			return err
		}
		recipe = recipeFound
		recipeIngredients, err := repo.RecipeViews().FindIngredients(ctx, id)
		if err != nil {
			return err
		}
		recipeIngredientsView = recipeIngredients
		return nil
	}); err != nil {
		return model.RecipeView{}, err
	}
	return model.RecipeView{
		ID:           recipe.ID,
		Name:         recipe.Name,
		Ingredients:  recipeIngredientsView,
		CreatedAt:    recipe.CreatedAt,
		LastModified: recipe.LastModified,
	}, nil
}
