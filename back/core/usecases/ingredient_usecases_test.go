package usecases_test

import (
	"context"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateIngredient(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should create an ingredient if non existent", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		clockMock := new(clockMock)
		now := time.UnixMilli(12345).UTC()
		clockMock.On("Now").Return(now)

		ingredientRepository := rpst.NewIngredientRepository(db, clockMock, logger)
		ingredientUsecases := usecases.NewIngredientUseCases(ingredientRepository, clockMock)
		ingredient, err := ingredientUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
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

	t.Run("should fail to create an ingredient if existent", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)
		ingredientUsecases := usecases.NewIngredientUseCases(ingredientRepository, clock)
		existentIngredientName := "name"
		ingredientUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
			Name:  existentIngredientName,
			Price: 10.0,
			Unit:  model.Gram,
		})

		_, err := ingredientUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
			Name:  existentIngredientName,
			Price: 1123450.0,
			Unit:  model.Kilogram,
		})
		require.Error(t, err)
	})

	t.Run("should assign different IDs to different ingredients", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)
		ingredientUsecases := usecases.NewIngredientUseCases(ingredientRepository, clock)
		ctx := context.Background()
		ing1, err := ingredientUsecases.CreateIngredient(ctx, usecases.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)
		ing2, err := ingredientUsecases.CreateIngredient(ctx, usecases.CreateIngredientOptions{
			Name:  "ing2",
			Price: 1231231231.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)

		assert.NotEqual(t, ing1.ID, ing2.ID)
	})

}

func TestEditIngredient(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()
	t.Run("should edit ingredient correctly", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)
		ingredientUsecases := usecases.NewIngredientUseCases(ingredientRepository, clock)
		ctx := context.Background()
		ing1, err := ingredientUsecases.CreateIngredient(ctx, usecases.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)

		newIngredientOpts := usecases.CreateIngredientOptions{
			Name:  "modifiedIngr1",
			Price: ing1.Price + 10.0,
			Unit:  model.Kilogram,
		}
		err = ingredientUsecases.EditIngredient(ctx, ing1.ID, newIngredientOpts)
		require.NoError(t, err)

		modifiedIngredient, err := ingredientRepository.GetIngredient(ctx, ing1.ID)
		require.NoError(t, err)

		assert.Equal(t, modifiedIngredient.Name, newIngredientOpts.Name)
		assert.Equal(t, modifiedIngredient.Price, newIngredientOpts.Price)
		assert.Equal(t, modifiedIngredient.Unit, newIngredientOpts.Unit)
	})
}
