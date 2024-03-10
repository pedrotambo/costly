package usecases

import (
	"context"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/rpst"
	"fmt"
)

type RecipeUseCases interface {
	RecipeCreator
}

type RecipeIngredientOptions struct {
	ID    int64
	Units int
}

type CreateRecipeOptions struct {
	Name        string
	Ingredients []RecipeIngredientOptions
}

type RecipeCreator interface {
	CreateRecipe(ctx context.Context, recipeOpts CreateRecipeOptions) (*domain.Recipe, error)
}

type recipeUseCases struct {
	repository rpst.Repository
	clock      clock.Clock
}

func NewRecipeUseCases(repository rpst.Repository, clock clock.Clock) RecipeUseCases {
	return &recipeUseCases{
		repository: repository,
		clock:      clock,
	}
}

func (cr *recipeUseCases) CreateRecipe(ctx context.Context, recipeOpts CreateRecipeOptions) (*domain.Recipe, error) {
	now := cr.clock.Now()
	recipeIngredients := []domain.RecipeIngredient{}
	for _, recipeIngredient := range recipeOpts.Ingredients {
		ingredient, err := cr.repository.GetIngredient(ctx, recipeIngredient.ID)
		if err == rpst.ErrNotFound {
			return &domain.Recipe{}, fmt.Errorf("failed to create recipe: unexistent ingredient with ID %d", recipeIngredient.ID)
		} else if err != nil {
			return &domain.Recipe{}, err
		}

		recipeIngredients = append(recipeIngredients, domain.RecipeIngredient{
			Ingredient: ingredient,
			Units:      recipeIngredient.Units,
		})
	}

	newRecipe := &domain.Recipe{
		ID:           -1,
		Name:         recipeOpts.Name,
		Ingredients:  recipeIngredients,
		CreatedAt:    now,
		LastModified: now,
	}

	err := cr.repository.SaveRecipe(ctx, newRecipe)

	if err != nil {
		return &domain.Recipe{}, fmt.Errorf("failed to create recipe: %s", err)
	}

	return newRecipe, nil
}
