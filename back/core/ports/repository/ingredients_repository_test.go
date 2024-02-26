package repository_test

import (
	"context"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type clockMock struct {
	mock.Mock
}

func (m *clockMock) Now() time.Time {
	args := m.Called()
	value := args.Get(0)
	now, ok := value.(time.Time)
	if !ok {
		panic(fmt.Errorf("Error getting now"))
	}
	return now
}

func TestIngredientRepository(t *testing.T) {

	logger, _ := logger.NewLogger("debug")

	t.Run("should create an ingredient if non existent", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		clockMock := new(clockMock)
		now := time.UnixMilli(12345).UTC()
		clockMock.On("Now").Return(now)

		ingredientRepository := repository.NewIngredientRepository(db, clockMock, logger)

		ingredient, err := ingredientRepository.CreateIngredient(context.Background(), "name", 10.0, domain.Gram)

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

		clock := clock.New()

		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)
		existentIngredientName := "name"
		ingredientRepository.CreateIngredient(context.Background(), existentIngredientName, 10.0, domain.Gram)

		_, err := ingredientRepository.CreateIngredient(context.Background(), existentIngredientName, 11231235.0, domain.Gram)
		require.Error(t, err)
	})

	t.Run("should get correct ingredient if existent", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		clock := clock.New()

		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)
		ctx := context.Background()
		ing1, err := ingredientRepository.CreateIngredient(ctx, "ing1", 10.0, domain.Gram)
		require.NoError(t, err)
		_, err = ingredientRepository.CreateIngredient(ctx, "ing2", 11231235.0, domain.Gram)
		require.NoError(t, err)

		ingr1Get, err := ingredientRepository.GetIngredient(ctx, ing1.ID)
		require.NoError(t, err)

		assert.Equal(t, ing1, ingr1Get)
	})

	t.Run("should assign different IDs to different ingredients", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		clock := clock.New()

		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)
		ctx := context.Background()
		ing1, err := ingredientRepository.CreateIngredient(ctx, "ing1", 10.0, domain.Gram)
		require.NoError(t, err)
		ing2, err := ingredientRepository.CreateIngredient(ctx, "ing2", 11231235.0, domain.Gram)
		require.NoError(t, err)

		assert.NotEqual(t, ing1.ID, ing2.ID)
	})

	t.Run("should return error when requesting an inexistent ingredient", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		clock := clock.New()
		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)

		_, err := ingredientRepository.GetIngredient(context.Background(), 123)

		require.Error(t, err)
		assert.Equal(t, err, repository.ErrNotFound)
	})

	t.Run("should edit ingredient price correctly", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		clock := clock.New()

		ingredientRepository := repository.NewIngredientRepository(db, clock, logger)
		ctx := context.Background()
		ing1, err := ingredientRepository.CreateIngredient(ctx, "ing1", 10.0, domain.Gram)
		require.NoError(t, err)

		newPrice := ing1.Price + 10.0
		newName := "modifiedIngr1"
		newUnit := domain.Kilogram
		ingredientRepository.EditIngredient(ctx, ing1.ID, newName, newPrice, newUnit)

		modifiedIngredient, err := ingredientRepository.GetIngredient(ctx, ing1.ID)
		require.NoError(t, err)

		assert.Equal(t, modifiedIngredient.Name, newName)
		assert.Equal(t, modifiedIngredient.Price, newPrice)
		assert.Equal(t, modifiedIngredient.Unit, newUnit)
	})
}
