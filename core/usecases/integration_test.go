package usecases_test

import (
	"context"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/usecases/ingredients"
	"costly/core/usecases/recipes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (ingredients.IngredientUseCases, recipes.RecipeUseCases, context.Context) {
	logger, err := logger.New("debug")
	require.NoError(t, err)
	clock := clock.New()
	db, err := database.NewFromDatasource(":memory:", logger)
	require.NoError(t, err)
	ingredients := ingredients.New(db, clock)
	recipes := recipes.New(db, clock, logger, ingredients)
	return ingredients, recipes, context.Background()
}

func TestIng1(t *testing.T) {
	ingredientUsecases, recipeUsecases, ctx := setupTest(t)

	ingredientUsecases.Create(ctx, ingredients.CreateIngredientOptions{
		Name:  "ing1",
		Price: 1.0,
		Unit:  model.Gram,
	})

	ingredientUsecases.Create(ctx, ingredients.CreateIngredientOptions{
		Name:  "ing2",
		Price: 2.0,
		Unit:  model.Gram,
	})

	ingrs, _ := ingredientUsecases.FindAll(ctx)

	assert.Equal(t, 0, ingrs[0].UnitsInStock)
	assert.Equal(t, 0, ingrs[1].UnitsInStock)

	recipe, err := recipeUsecases.Create(ctx, recipes.CreateRecipeOptions{
		Name: "recipe1",
		Ingredients: []model.RecipeIngredient{
			{
				ID:    ingrs[0].ID,
				Units: 2,
			},
			{
				ID:    ingrs[1].ID,
				Units: 5,
			},
		},
	})

	require.NoError(t, err)

	_, err = recipeUsecases.AddSales(ctx, recipe.ID, 2)
	require.NoError(t, err)

	ingrs, _ = ingredientUsecases.FindAll(ctx)
	assert.Equal(t, -4, ingrs[0].UnitsInStock)
	assert.Equal(t, -10, ingrs[1].UnitsInStock)

	ingredientUsecases.AddStock(ctx, ingrs[0].ID, ingredients.IngredientStockOptions{
		Units: 4,
		Price: 1.0,
	})

	ingredientUsecases.AddStock(ctx, ingrs[1].ID, ingredients.IngredientStockOptions{
		Units: 10,
		Price: 2.0,
	})

	ingrs, _ = ingredientUsecases.FindAll(ctx)
	assert.Equal(t, 0, ingrs[0].UnitsInStock)
	assert.Equal(t, 0, ingrs[1].UnitsInStock)

	fmt.Println(ingrs)
}
