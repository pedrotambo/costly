package recipes_test

import (
	"context"
	"costly/core/errs"
	"costly/core/mocks"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/usecases/ingredients"
	"costly/core/usecases/recipes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var meat = ingredients.CreateIngredientOptions{
	Name:  "meat",
	Price: 1.0,
	Unit:  model.Gram,
}

var salt = ingredients.CreateIngredientOptions{
	Name:  "salt",
	Price: 10.0,
	Unit:  model.Gram,
}

var pepper = ingredients.CreateIngredientOptions{
	Name:  "pepper",
	Price: 13.0,
	Unit:  model.Gram,
}

func setupTest(logger logger.Logger, clock clock.Clock) ([]model.Ingredient, recipes.RecipeUseCases, context.Context) {
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientUseCases := ingredients.New(db, clock)
	ctx := context.Background()
	var createdIngredients = []model.Ingredient{}
	for _, ingredient := range []ingredients.CreateIngredientOptions{meat, salt, pepper} {
		ing, _ := ingredientUseCases.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  ingredient.Name,
			Price: ingredient.Price,
			Unit:  ingredient.Unit,
		})
		createdIngredients = append(createdIngredients, *ing)
	}
	recipeUseCases := recipes.New(db, clock, logger, ingredientUseCases)
	return createdIngredients, recipeUseCases, context.Background()
}

func TestCreate(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should create a recipe if non existent", func(t *testing.T) {
		clockMock := new(mocks.ClockMock)
		now := time.UnixMilli(12345).UTC()
		clockMock.On("Now").Return(now)
		ingredients, recipeComponent, ctx := setupTest(logger, clockMock)

		recipe, err := recipeComponent.Create(ctx, recipes.CreateRecipeOptions{
			Name: "recipeName",
			Ingredients: []model.RecipeIngredient{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				}, {
					ID:    ingredients[2].ID,
					Units: 5,
				},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), recipe.ID)
		assert.Equal(t, "recipeName", recipe.Name)
		assert.Len(t, recipe.Ingredients, 2)
		for _, recipeIngredient := range recipe.Ingredients {
			if recipeIngredient.ID == ingredients[0].ID {
				assert.Equal(t, 500, recipeIngredient.Units)
			} else {
				assert.Equal(t, 5, recipeIngredient.Units)
			}

		}
		assert.Equal(t, recipe.CreatedAt, now)
		assert.Equal(t, recipe.LastModified, now)
	})

	t.Run("should fail to create a recipe if existent", func(t *testing.T) {
		ingredients, recipeComponent, ctx := setupTest(logger, clock)
		existentRecipeName := "name"

		recipeComponent.Create(ctx, recipes.CreateRecipeOptions{
			Name: existentRecipeName,
			Ingredients: []model.RecipeIngredient{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				},
			},
		})

		_, err := recipeComponent.Create(ctx, recipes.CreateRecipeOptions{
			Name: existentRecipeName,
			Ingredients: []model.RecipeIngredient{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				},
			},
		})
		require.Error(t, err)
		assert.EqualError(t, err, "failed to create recipe: UNIQUE constraint failed: recipe.name")
	})

	t.Run("should return an error when creating a recipe with unexistent ingredient", func(t *testing.T) {
		ingredients, recipeComponent, ctx := setupTest(logger, clock)
		existentRecipeName := "name"
		var unexistentIngredientID int64
		for _, i := range ingredients {
			unexistentIngredientID += i.ID
		}
		_, err := recipeComponent.Create(ctx, recipes.CreateRecipeOptions{
			Name: existentRecipeName,
			Ingredients: []model.RecipeIngredient{
				{
					ID:    unexistentIngredientID,
					Units: 500,
				},
			},
		})
		assert.EqualError(t, err, "failed to create recipe: FOREIGN KEY constraint failed")
	})

	t.Run("should assign different IDs to different recipes", func(t *testing.T) {
		ingredients, recipeComponent, ctx := setupTest(logger, clock)
		recipe1, err := recipeComponent.Create(ctx, recipes.CreateRecipeOptions{
			Name: "recipe1",
			Ingredients: []model.RecipeIngredient{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				},
			},
		})
		require.NoError(t, err)
		recipe2, err := recipeComponent.Create(ctx, recipes.CreateRecipeOptions{
			Name: "recipe2",
			Ingredients: []model.RecipeIngredient{
				{
					ID:    ingredients[1].ID,
					Units: 200,
				},
			},
		})
		require.NoError(t, err)

		assert.NotEqual(t, recipe1.ID, recipe2.ID)
	})

	t.Run("should return error when creating a recipe without ingredients", func(t *testing.T) {
		_, recipeComponent, ctx := setupTest(logger, clock)
		_, err := recipeComponent.Create(ctx, recipes.CreateRecipeOptions{
			Name:        "recipe1",
			Ingredients: []model.RecipeIngredient{},
		})
		require.Error(t, err)
		assert.Equal(t, errs.ErrBadIngrs, err)
	})
}
