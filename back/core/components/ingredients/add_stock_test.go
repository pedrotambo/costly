package ingredients_test

import (
	"costly/core/components/ingredients"
	"costly/core/errs"
	"costly/core/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddStock(t *testing.T) {

	t.Run("should return error if bad price", func(t *testing.T) {
		ingredientComponent, ctx := setupTest(t)
		_, err := ingredientComponent.AddStock(ctx, 1, ingredients.IngredientStockOptions{
			Price: 0.0,
			Units: 5,
		})
		assert.Equal(t, errs.ErrBadPrice, err)
	})

	t.Run("should return error if bad units", func(t *testing.T) {
		ingredientComponent, ctx := setupTest(t)
		_, err := ingredientComponent.AddStock(ctx, 1, ingredients.IngredientStockOptions{
			Price: 1.0,
			Units: 0,
		})
		assert.Equal(t, errs.ErrBadStockUnits, err)
	})

	t.Run("should update price and units in stock", func(t *testing.T) {
		ingredientComponent, ctx := setupTest(t)
		ingredient, err := ingredientComponent.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)

		ingredientComponent.AddStock(ctx, ingredient.ID, ingredients.IngredientStockOptions{
			Price: 1.0,
			Units: 5,
		})
		ingredientComponent.AddStock(ctx, ingredient.ID, ingredients.IngredientStockOptions{
			Price: 2.0,
			Units: 7,
		})

		modifiedIngredient, err := ingredientComponent.Find(ctx, ingredient.ID)
		require.NoError(t, err)
		assert.Equal(t, ingredient.UnitsInStock+5+7, modifiedIngredient.UnitsInStock)
		assert.Equal(t, 2.0, modifiedIngredient.Price)
	})

	t.Run("should add new stock in history", func(t *testing.T) {
		ingredientComponent, ctx := setupTest(t)
		ingredient, err := ingredientComponent.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)

		stock, err := ingredientComponent.AddStock(ctx, ingredient.ID, ingredients.IngredientStockOptions{
			Price: 1.0,
			Units: 5,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), stock.ID)
		assert.Equal(t, ingredient.ID, stock.IngredientID)
		assert.Equal(t, 1.0, stock.Price)
		assert.Equal(t, 5, stock.Units)
	})

	t.Run("should return error if unexistent ingredient", func(t *testing.T) {
		ingredientComponent, ctx := setupTest(t)
		_, err := ingredientComponent.AddStock(ctx, 1, ingredients.IngredientStockOptions{
			Price: 1.0,
			Units: 5,
		})
		require.Error(t, err)
		assert.Equal(t, errs.ErrNotFound, err)
	})
}
