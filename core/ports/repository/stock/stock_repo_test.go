package stockrepo_test

import (
	"context"

	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	ingredientrepo "costly/core/ports/repository/ingredient"
	stockrepo "costly/core/ports/repository/stock"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddIngredientStock(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should return error if unexistent ingredient", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		stockRepository := stockrepo.New(db)
		ctx := context.Background()
		stock, _ := model.NewIngredientStock(1, 5, 1.0, clock.Now())
		err := stockRepository.Add(ctx, stock)
		require.Error(t, err)
	})

	t.Run("should add findable entity", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		ingredientRepository := ingredientrepo.New(db)
		stockRepository := stockrepo.New(db)
		ctx := context.Background()
		now := clock.Now()

		ingredient, _ := model.NewIngredient("first", model.Gram, 1.0, now)
		err := ingredientRepository.Add(ctx, ingredient)
		require.NoError(t, err)
		ingredient2, _ := model.NewIngredient("second", model.Gram, 1.0, now)
		err = ingredientRepository.Add(ctx, ingredient2)
		require.NoError(t, err)
		stock, _ := model.NewIngredientStock(ingredient.ID, 5, 1.0, clock.Now())
		stockRepository.Add(ctx, stock)
		stockGet, err := stockRepository.Find(ctx, stock.ID)
		require.NoError(t, err)

		assert.Equal(t, stock.ID, stockGet.ID)
		assert.Equal(t, ingredient.ID, stockGet.IngredientID)
		assert.Equal(t, 1.0, stockGet.Price)
		assert.Equal(t, 5, stockGet.Units)
	})
}
