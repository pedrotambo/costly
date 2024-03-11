package rpst_test

import (
	"context"
	"costly/core/components/ingredients/internal/rpst"
	"costly/core/errs"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
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

		ingredientRepository := rpst.New(db, logger)
		ctx := context.Background()
		now := clock.Now()
		ingredient := &model.Ingredient{
			ID:           1,
			Name:         "aName",
			Unit:         model.Gram,
			Price:        1.0,
			UnitsInStock: 0,
			CreatedAt:    now,
			LastModified: now,
		}
		err := ingredientRepository.Add(ctx, ingredient)
		require.NoError(t, err)

		ingredient2 := &model.Ingredient{
			ID:           1,
			Name:         "ing2",
			Unit:         model.Gram,
			Price:        1123123123120.0,
			UnitsInStock: 0,
			CreatedAt:    now,
			LastModified: now,
		}
		err = ingredientRepository.Add(ctx, ingredient2)
		require.NoError(t, err)

		ingr1Get, err := ingredientRepository.Find(ctx, ingredient.ID)
		require.NoError(t, err)

		assert.Equal(t, ingredient, &ingr1Get)
	})

	t.Run("should return error when requesting an inexistent ingredient", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		ingredientRepository := rpst.New(db, logger)

		_, err := ingredientRepository.Find(context.Background(), 123)

		require.Error(t, err)
		assert.Equal(t, err, errs.ErrNotFound)
	})
}

func TestGetIngredients(t *testing.T) {

	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get list of existent ingredients", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := rpst.New(db, logger)
		ctx := context.Background()
		now := clock.Now()
		ingredient := &model.Ingredient{
			ID:           1,
			Name:         "aName",
			Unit:         model.Gram,
			Price:        1.0,
			UnitsInStock: 0,
			CreatedAt:    now,
			LastModified: now,
		}
		err := ingredientRepository.Add(ctx, ingredient)
		require.NoError(t, err)

		ingredient2 := &model.Ingredient{
			ID:           1,
			Name:         "ing2",
			Unit:         model.Gram,
			Price:        1123123123120.0,
			UnitsInStock: 0,
			CreatedAt:    now,
			LastModified: now,
		}
		err = ingredientRepository.Add(ctx, ingredient2)
		require.NoError(t, err)

		ingredients, err := ingredientRepository.FindAll(ctx)
		require.NoError(t, err)

		assert.Equal(t, ingredient, &ingredients[0])
		assert.Equal(t, ingredient2, &ingredients[1])
	})
}

func TestAddIngredientStock(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should add ingredient units in stock correctly", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := rpst.New(db, logger)
		ctx := context.Background()
		now := clock.Now()
		ingredient := &model.Ingredient{
			ID:           1,
			Name:         "aName",
			Unit:         model.Gram,
			Price:        1.0,
			UnitsInStock: 0,
			CreatedAt:    now,
			LastModified: now,
		}
		err := ingredientRepository.Add(ctx, ingredient)
		require.NoError(t, err)

		ingredientRepository.AddStock(ctx, &model.IngredientStock{
			ID:           -1,
			IngredientID: ingredient.ID,
			Price:        1.0,
			Units:        5,
			CreatedAt:    clock.Now(),
		})
		ingredientRepository.AddStock(ctx, &model.IngredientStock{
			ID:           -1,
			IngredientID: ingredient.ID,
			Price:        2.0,
			Units:        7,
			CreatedAt:    clock.Now(),
		})

		modifiedIngredient, err := ingredientRepository.Find(ctx, ingredient.ID)
		require.NoError(t, err)

		assert.Equal(t, ingredient.UnitsInStock+5+7, modifiedIngredient.UnitsInStock)
	})

	t.Run("when adding ingredient stock should update price correctly", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)

		ingredientRepository := rpst.New(db, logger)
		ctx := context.Background()
		now := clock.Now()
		ingredient := &model.Ingredient{
			ID:           1,
			Name:         "aName",
			Unit:         model.Gram,
			Price:        1.0,
			UnitsInStock: 0,
			CreatedAt:    now,
			LastModified: now,
		}
		err := ingredientRepository.Add(ctx, ingredient)
		require.NoError(t, err)

		ingredientRepository.AddStock(ctx, &model.IngredientStock{
			ID:           -1,
			IngredientID: ingredient.ID,
			Price:        1.0,
			Units:        5,
			CreatedAt:    clock.Now(),
		})
		ingredientRepository.AddStock(ctx, &model.IngredientStock{
			ID:           -1,
			IngredientID: ingredient.ID,
			Price:        2.0,
			Units:        7,
			CreatedAt:    clock.Now(),
		})

		modifiedIngredient, err := ingredientRepository.Find(ctx, ingredient.ID)
		require.NoError(t, err)

		assert.Equal(t, 2.0, modifiedIngredient.Price)
	})

	t.Run("adding stock of inexistent ingredient should return error", func(t *testing.T) {
		db, _ := database.NewFromDatasource(":memory:", logger)
		ingredientRepository := rpst.New(db, logger)
		ctx := context.Background()
		now := clock.Now()
		ingredient := &model.Ingredient{
			ID:           1,
			Name:         "aName",
			Unit:         model.Gram,
			Price:        1.0,
			UnitsInStock: 0,
			CreatedAt:    now,
			LastModified: now,
		}
		err := ingredientRepository.Add(ctx, ingredient)
		require.NoError(t, err)

		ingredientRepository.AddStock(ctx, &model.IngredientStock{
			ID:           -1,
			IngredientID: ingredient.ID,
			Price:        1.0,
			Units:        5,
			CreatedAt:    clock.Now(),
		})
		require.NoError(t, err)

		err = ingredientRepository.AddStock(ctx, &model.IngredientStock{
			ID:           -1,
			IngredientID: ingredient.ID + 1,
			Price:        1.0,
			Units:        5,
			CreatedAt:    clock.Now(),
		})
		require.Error(t, err)
		assert.Equal(t, err, errs.ErrNotFound)
	})

	t.Run("add ingredient stock should return error if query returns error", func(t *testing.T) {
		db := new(databaseMock)
		ingredientRepository := rpst.New(db, logger)
		err := ingredientRepository.AddStock(context.Background(), &model.IngredientStock{
			ID:           -1,
			IngredientID: 1,
			Price:        1.0,
			Units:        5,
			CreatedAt:    clock.Now(),
		})
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
