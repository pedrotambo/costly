package ingredients_test

import (
	"context"
	"costly/core/errs"
	"costly/core/mocks"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/usecases/ingredients"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (ingredients.IngredientUseCases, context.Context) {
	logger, err := logger.New("debug")
	require.NoError(t, err)
	clock := clock.New()
	db, err := database.NewFromDatasource(":memory:", logger)
	require.NoError(t, err)
	return ingredients.New(db, clock), context.Background()
}

func TestCreateIngredient(t *testing.T) {

	t.Run("should create an ingredient if non existent", func(t *testing.T) {
		logger, _ := logger.New("debug")
		db, _ := database.NewFromDatasource(":memory:", logger)
		clockMock := new(mocks.ClockMock)
		now := time.UnixMilli(12345).UTC()
		clockMock.On("Now").Return(now)
		ingredientUseCases := ingredients.New(db, clockMock)
		ingredient, err := ingredientUseCases.Create(context.Background(), ingredients.CreateIngredientOptions{
			Name:  "name",
			Price: 10.0,
			Unit:  model.Gram,
		})

		if err != nil {
			t.Fatal()
		}
		assert.Equal(t, ingredient.ID, int64(1))
		assert.Equal(t, ingredient.Name, "name")
		assert.Equal(t, ingredient.Price, 10.0)
		assert.Equal(t, ingredient.Unit, model.Gram)
		assert.Equal(t, ingredient.CreatedAt, now)
		assert.Equal(t, ingredient.LastModified, now)
	})

	t.Run("should return error if invalid options", func(t *testing.T) {
		ingredientComponent, ctx := setupTest(t)
		_, err := ingredientComponent.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  "",
			Price: 10.0,
			Unit:  model.Gram,
		})
		assert.Equal(t, err, errs.ErrBadName)
	})

	t.Run("should fail to create an ingredient if existent", func(t *testing.T) {
		ingredientComponent, ctx := setupTest(t)
		existentIngredientName := "name"
		ingredientComponent.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  existentIngredientName,
			Price: 10.0,
			Unit:  model.Gram,
		})

		_, err := ingredientComponent.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  existentIngredientName,
			Price: 1123450.0,
			Unit:  model.Kilogram,
		})
		require.Error(t, err)
	})

	t.Run("should assign different IDs to different ingredients", func(t *testing.T) {
		ingredientComponent, ctx := setupTest(t)
		ing1, err := ingredientComponent.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)
		ing2, err := ingredientComponent.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  "ing2",
			Price: 1231231231.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)

		assert.NotEqual(t, ing1.ID, ing2.ID)
	})
}
