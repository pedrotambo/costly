package model_test

import (
	"costly/core/errs"
	"costly/core/model"
	"costly/core/ports/clock"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipeCost(t *testing.T) {

	t.Run("cost of a recipe is the sum of its ingredients and units of them", func(t *testing.T) {
		clock := clock.New()
		now := clock.Now()
		recipe := model.RecipeView{
			ID:   1,
			Name: "aName",
			Ingredients: []model.RecipeIngredientView{
				{
					ID:    1,
					Name:  "meat",
					Units: 500,
					Price: 1,
				},
				{
					ID:    2,
					Name:  "salt",
					Units: 5,
					Price: 10.0,
				},
			},
			CreatedAt:    now,
			LastModified: now,
		}

		assert.Equal(t, recipe.Cost(), 500.0+5*10)
	})
}

func TestNewIngredient(t *testing.T) {

	now := clock.New().Now()
	t.Run("should add invalid ID and the rest of the fields correct", func(t *testing.T) {
		ingredient, err := model.NewIngredient("name", model.Gram, 123.0, now)
		require.NoError(t, err)
		assert.Equal(t, int64(-1), ingredient.ID)
		assert.Equal(t, "name", ingredient.Name)
		assert.Equal(t, model.Gram, ingredient.Unit)
		assert.Equal(t, 123.0, ingredient.Price)
		assert.Equal(t, now, ingredient.CreatedAt)
	})

	t.Run("should return error if name is invalid", func(t *testing.T) {
		_, err := model.NewIngredient("", model.Gram, 123.0, now)
		assert.Equal(t, err, errs.ErrBadName)
	})

	t.Run("should return error if unit is invalid", func(t *testing.T) {
		_, err := model.NewIngredient("name", "asdf", 123.0, now)
		assert.Equal(t, err, errs.ErrBadUnit)
	})

	t.Run("should return error if price is invalid", func(t *testing.T) {
		_, err := model.NewIngredient("name", model.Gram, 0, now)
		assert.Equal(t, err, errs.ErrBadPrice)
	})
}

func TestNewIngredientStock(t *testing.T) {

	now := clock.New().Now()
	t.Run("should add invalid ID and the rest of the fields correct", func(t *testing.T) {
		stock, err := model.NewIngredientStock(int64(1), 1, 123.0, now)
		require.NoError(t, err)
		assert.Equal(t, int64(-1), stock.ID)
		assert.Equal(t, int64(1), stock.IngredientID)
		assert.Equal(t, 1, stock.Units)
		assert.Equal(t, 123.0, stock.Price)
		assert.Equal(t, now, stock.CreatedAt)
	})

	t.Run("should return error if units is invalid", func(t *testing.T) {
		_, err := model.NewIngredientStock(int64(1), 0, 123.0, now)
		assert.Equal(t, err, errs.ErrBadStockUnits)
	})

	t.Run("should return error if price is invalid", func(t *testing.T) {
		_, err := model.NewIngredientStock(int64(1), 1, 0.0, now)
		assert.Equal(t, err, errs.ErrBadPrice)
	})
}

func TestNewRecipeSales(t *testing.T) {

	now := clock.New().Now()
	t.Run("should add invalid ID and the rest of the fields correct", func(t *testing.T) {
		sales := model.NewRecipeSales(int64(1), 1, now)
		assert.Equal(t, int64(-1), sales.ID)
		assert.Equal(t, int64(1), sales.RecipeID)
		assert.Equal(t, 1, sales.Units)
		assert.Equal(t, now, sales.CreatedAt)
	})
}

func TestNewRecipe(t *testing.T) {

	now := clock.New().Now()
	t.Run("should add invalid ID and the rest of the fields correct", func(t *testing.T) {
		recipeIngredients := []model.RecipeIngredient{
			{
				ID:    5,
				Units: 150,
			},
		}
		recipe, err := model.NewRecipe("name", recipeIngredients, now)
		require.NoError(t, err)
		assert.Equal(t, int64(-1), recipe.ID)
		assert.Equal(t, "name", recipe.Name)
		assert.Equal(t, recipeIngredients, recipe.Ingredients)
		assert.Equal(t, now, recipe.CreatedAt)
		assert.Equal(t, now, recipe.LastModified)
	})

	t.Run("should return error if name is invalid", func(t *testing.T) {
		recipeIngredients := []model.RecipeIngredient{
			{
				ID:    5,
				Units: 150,
			},
		}
		_, err := model.NewRecipe("", recipeIngredients, now)
		assert.Equal(t, err, errs.ErrBadName)
	})

	t.Run("should return error if ingredients is empty", func(t *testing.T) {
		_, err := model.NewRecipe("name", []model.RecipeIngredient{}, now)
		assert.Equal(t, err, errs.ErrBadIngrs)
	})
}
