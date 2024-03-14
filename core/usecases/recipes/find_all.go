package recipes

import (
	"context"
	"costly/core/model"
)

type RecipesFinder interface {
	FindAll(ctx context.Context) ([]model.RecipeView, error)
}

func (cr *recipeUseCases) FindAll(ctx context.Context) ([]model.RecipeView, error) {
	return cr.repository.RecipeViews().FindAll(ctx)
}
