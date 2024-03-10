package rpst_test

import (
	"context"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetIngredient(t *testing.T) {

	logger, _ := logger.New("debug")
	clock := clock.New()
	t.Run("should get correct ingredient if existent", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)
		ingredientUsecases := usecases.NewIngredientUseCases(ingredientRepository, clock)
		ctx := context.Background()
		ing1, err := ingredientUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)
		_, err = ingredientUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
			Name:  "ing2",
			Price: 1123123123120.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)

		ingr1Get, err := ingredientRepository.GetIngredient(ctx, ing1.ID)
		require.NoError(t, err)

		assert.Equal(t, ing1, &ingr1Get)
	})

	t.Run("should return error when requesting an inexistent ingredient", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)

		_, err := ingredientRepository.GetIngredient(context.Background(), 123)

		require.Error(t, err)
		assert.Equal(t, err, rpst.ErrNotFound)
	})
}

func TestGetIngredients(t *testing.T) {

	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get list of existent ingredients", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)
		ingredientUsecases := usecases.NewIngredientUseCases(ingredientRepository, clock)
		ctx := context.Background()
		ing1, err := ingredientUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)
		ing2, err := ingredientUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
			Name:  "ing2",
			Price: 1123123123120.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)

		ingredients, err := ingredientRepository.GetIngredients(ctx)
		require.NoError(t, err)

		assert.Equal(t, ing1, &ingredients[0])
		assert.Equal(t, ing2, &ingredients[1])
	})
}

func TestUpdateStock(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should update ingredient units in stock correctly", func(t *testing.T) {
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

		ingredientRepository.UpdateStock(ctx, ing1.ID, rpst.NewStockOptions{NewUnits: 5, Price: 1.0})
		ingredientRepository.UpdateStock(ctx, ing1.ID, rpst.NewStockOptions{NewUnits: 7, Price: 2.0})

		modifiedIngredient, err := ingredientRepository.GetIngredient(ctx, ing1.ID)
		require.NoError(t, err)

		assert.Equal(t, ing1.UnitsInStock+5+7, modifiedIngredient.UnitsInStock)
	})

	t.Run("should update ingredient price correctly", func(t *testing.T) {
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

		ingredientRepository.UpdateStock(ctx, ing1.ID, rpst.NewStockOptions{NewUnits: 5, Price: 1.0})
		ingredientRepository.UpdateStock(ctx, ing1.ID, rpst.NewStockOptions{NewUnits: 7, Price: 2.0})

		modifiedIngredient, err := ingredientRepository.GetIngredient(ctx, ing1.ID)
		require.NoError(t, err)

		assert.Equal(t, 2.0, modifiedIngredient.Price)
	})

	t.Run("update stock of inexistent ingredient should return error", func(t *testing.T) {
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

		ingredientRepository.UpdateStock(ctx, ing1.ID, rpst.NewStockOptions{NewUnits: 5, Price: 1.0})
		require.NoError(t, err)

		_, err = ingredientRepository.UpdateStock(ctx, ing1.ID+1, rpst.NewStockOptions{NewUnits: 7, Price: 1.0})
		require.Error(t, err)
		assert.Equal(t, err, rpst.ErrNotFound)
	})

	t.Run("update stock should return error if query returns error", func(t *testing.T) {
		db := new(databaseMock)
		ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)
		ctx := context.Background()

		_, err := ingredientRepository.UpdateStock(ctx, 1, rpst.NewStockOptions{NewUnits: 7, Price: 1.0})
		require.Error(t, err)
		assert.Equal(t, err, ErrDBInternal)
	})
}

var ErrDBInternal = errors.New("internal db error")

type databaseMock struct {
	mock.Mock
}

type errorRow struct{}

func (e *errorRow) Scan(dest ...any) error {
	return ErrDBInternal
}

func (e *errorRow) Next() bool {
	return false
}

func (dm *databaseMock) QueryRowContext(ctx context.Context, query string, args ...any) database.RowScanner {
	return &errorRow{}
}

func (dm *databaseMock) QueryContext(ctx context.Context, query string, args ...any) (database.RowsScanner, error) {
	return nil, ErrDBInternal
}

func (dm *databaseMock) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return nil, ErrDBInternal
}

func (dm *databaseMock) WithTx(ctx context.Context, op func(tx database.TX) error) error {
	return ErrDBInternal
}
