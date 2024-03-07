package repository_test

import (
	"context"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIngredientRepository(t *testing.T) {

	logger, _ := logger.NewLogger("debug")
	clock := clock.New()

	t.Run("should create an ingredient if non existent", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		clockMock := new(clockMock)
		now := time.UnixMilli(12345).UTC()
		clockMock.On("Now").Return(now)

		ingredientRepository := repository.NewIngredientRepository(db, clockMock, logger)

		ingredient, err := ingredientRepository.CreateIngredient(context.Background(), repository.CreateIngredientOptions{
			Name:  "name",
			Price: 10.0,
			Unit:  domain.Gram,
		})

		if err != nil {
			t.Fail()
		}

		assert.Equal(t, ingredient.ID, int64(1))
		assert.Equal(t, ingredient.Name, "name")
		assert.Equal(t, ingredient.Price, 10.0)
		assert.Equal(t, ingredient.Unit, domain.Gram)
		assert.Equal(t, ingredient.CreatedAt, now)
		assert.Equal(t, ingredient.LastModified, now)
	})

	t.Run("should fail to create an ingredient if existent", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)
		existentIngredientName := "name"
		ingredientRepository.CreateIngredient(context.Background(), repository.CreateIngredientOptions{
			Name:  existentIngredientName,
			Price: 10.0,
			Unit:  domain.Gram,
		})

		_, err := ingredientRepository.CreateIngredient(context.Background(), repository.CreateIngredientOptions{
			Name:  existentIngredientName,
			Price: 1123450.0,
			Unit:  domain.Kilogram,
		})
		require.Error(t, err)
	})

	t.Run("should get correct ingredient if existent", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)
		ctx := context.Background()
		ing1, err := ingredientRepository.CreateIngredient(context.Background(), repository.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  domain.Gram,
		})
		require.NoError(t, err)
		_, err = ingredientRepository.CreateIngredient(context.Background(), repository.CreateIngredientOptions{
			Name:  "ing2",
			Price: 1123123123120.0,
			Unit:  domain.Gram,
		})
		require.NoError(t, err)

		ingr1Get, err := ingredientRepository.GetIngredient(ctx, ing1.ID)
		require.NoError(t, err)

		assert.Equal(t, ing1, ingr1Get)
	})

	t.Run("should assign different IDs to different ingredients", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)
		ctx := context.Background()
		ing1, err := ingredientRepository.CreateIngredient(ctx, repository.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  domain.Gram,
		})
		require.NoError(t, err)
		ing2, err := ingredientRepository.CreateIngredient(ctx, repository.CreateIngredientOptions{
			Name:  "ing2",
			Price: 1231231231.0,
			Unit:  domain.Gram,
		})
		require.NoError(t, err)

		assert.NotEqual(t, ing1.ID, ing2.ID)
	})

	t.Run("should return error when requesting an inexistent ingredient", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)

		_, err := ingredientRepository.GetIngredient(context.Background(), 123)

		require.Error(t, err)
		assert.Equal(t, err, repository.ErrNotFound)
	})

	t.Run("should edit ingredient correctly", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)
		ctx := context.Background()
		ing1, err := ingredientRepository.CreateIngredient(ctx, repository.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  domain.Gram,
		})
		require.NoError(t, err)

		newIngredientOpts := repository.CreateIngredientOptions{
			Name:  "modifiedIngr1",
			Price: ing1.Price + 10.0,
			Unit:  domain.Kilogram,
		}
		ingredientRepository.EditIngredient(ctx, ing1.ID, newIngredientOpts)

		modifiedIngredient, err := ingredientRepository.GetIngredient(ctx, ing1.ID)
		require.NoError(t, err)

		assert.Equal(t, modifiedIngredient.Name, newIngredientOpts.Name)
		assert.Equal(t, modifiedIngredient.Price, newIngredientOpts.Price)
		assert.Equal(t, modifiedIngredient.Unit, newIngredientOpts.Unit)
	})
}
