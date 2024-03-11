package usecases

import (
	"context"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/rpst"
	"fmt"
)

type RecipeUseCases interface {
	RecipeCreator
	RecipeSalesAdder
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
	CreateRecipe(ctx context.Context, recipeOpts CreateRecipeOptions) (*model.Recipe, error)
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

func (cr *recipeUseCases) CreateRecipe(ctx context.Context, recipeOpts CreateRecipeOptions) (*model.Recipe, error) {
	now := cr.clock.Now()
	recipeIngredients := []model.RecipeIngredient{}
	for _, recipeIngredient := range recipeOpts.Ingredients {
		ingredient, err := cr.repository.GetIngredient(ctx, recipeIngredient.ID)
		if err == rpst.ErrNotFound {
			return &model.Recipe{}, fmt.Errorf("failed to create recipe: unexistent ingredient with ID %d", recipeIngredient.ID)
		} else if err != nil {
			return &model.Recipe{}, err
		}

		recipeIngredients = append(recipeIngredients, model.RecipeIngredient{
			Ingredient: ingredient,
			Units:      recipeIngredient.Units,
		})
	}

	newRecipe := &model.Recipe{
		ID:           -1,
		Name:         recipeOpts.Name,
		Ingredients:  recipeIngredients,
		CreatedAt:    now,
		LastModified: now,
	}

	err := cr.repository.SaveRecipe(ctx, newRecipe)

	if err != nil {
		return &model.Recipe{}, fmt.Errorf("failed to create recipe: %s", err)
	}

	return newRecipe, nil
}

type RecipeSalesOpts struct {
	RecipeID int64
	Units    int
}

type RecipeSalesAdder interface {
	AddRecipeSales(ctx context.Context, recipeID int64, soldUnits int) (*model.RecipeSales, error)
}

func (cr *recipeUseCases) AddRecipeSales(ctx context.Context, recipeID int64, soldUnits int) (*model.RecipeSales, error) {
	now := cr.clock.Now()
	recipeSales := &model.RecipeSales{
		ID:        -1,
		RecipeID:  recipeID,
		Units:     soldUnits,
		CreatedAt: now,
	}
	err := cr.repository.SaveRecipeSales(ctx, recipeSales)

	if err != nil {
		return &model.RecipeSales{}, fmt.Errorf("failed to add recipe sales: %s", err)
	}

	return recipeSales, nil
}
