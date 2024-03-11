package recipes

import (
	"context"
	"costly/core/components/clock"
	"costly/core/components/database"
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/components/recipes/internal/rpst"
	"costly/core/errs"
	"costly/core/model"
	"fmt"
)

type RecipeComponent interface {
	RecipeCreator
	RecipeSalesAdder
	RecipeFinder
	RecipesFinder
}

type RecipeIngredientOptions struct {
	ID    int64
	Units int
}

type recipeUseCases struct {
	repository  rpst.RecipeRepository
	clock       clock.Clock
	ingredients ingredients.IngredientComponent
	rpst.RecipeRepository
}

func New(database database.Database, clock clock.Clock, logger logger.Logger, ingredients ingredients.IngredientComponent) RecipeComponent {
	recipeRepository := rpst.New(database, logger)
	return &recipeUseCases{
		repository:  recipeRepository,
		clock:       clock,
		ingredients: ingredients,
	}
}

type RecipeCreator interface {
	Create(ctx context.Context, recipeOpts CreateRecipeOptions) (*model.Recipe, error)
}

type CreateRecipeOptions struct {
	Name        string
	Ingredients []RecipeIngredientOptions
}

func (cr *recipeUseCases) Create(ctx context.Context, recipeOpts CreateRecipeOptions) (*model.Recipe, error) {
	now := cr.clock.Now()
	recipeIngredients := []model.RecipeIngredient{}
	for _, recipeIngredient := range recipeOpts.Ingredients {
		ingredient, err := cr.ingredients.Find(ctx, recipeIngredient.ID)
		if err == errs.ErrNotFound {
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

	err := cr.repository.Add(ctx, newRecipe)

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
	AddSales(ctx context.Context, recipeID int64, soldUnits int) (*model.RecipeSales, error)
}

func (cr *recipeUseCases) AddSales(ctx context.Context, recipeID int64, soldUnits int) (*model.RecipeSales, error) {
	now := cr.clock.Now()
	recipeSales := &model.RecipeSales{
		ID:        -1,
		RecipeID:  recipeID,
		Units:     soldUnits,
		CreatedAt: now,
	}
	err := cr.repository.AddSales(ctx, recipeSales)

	if err != nil {
		return &model.RecipeSales{}, fmt.Errorf("failed to add recipe sales: %s", err)
	}

	return recipeSales, nil
}

type RecipeFinder interface {
	Find(ctx context.Context, id int64) (model.Recipe, error)
}

func (cr *recipeUseCases) Find(ctx context.Context, id int64) (model.Recipe, error) {
	return cr.repository.Find(ctx, id)
}

type RecipesFinder interface {
	FindAll(ctx context.Context) ([]model.Recipe, error)
}

func (cr *recipeUseCases) FindAll(ctx context.Context) ([]model.Recipe, error) {
	return cr.repository.FindAll(ctx)
}
