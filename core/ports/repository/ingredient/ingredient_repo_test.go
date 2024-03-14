package ingredientrepo_test

import (
	"context"

	"costly/core/errs"
	"costly/core/mocks"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	ingredientrepo "costly/core/ports/repository/ingredient"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (ingredientrepo.IngredientRepository, clock.Clock, context.Context) {
	logger, _ := logger.New("debug")
	clock := clock.New()
	db, err := database.NewFromDatasource(":memory:", logger)
	if err != nil {
		t.FailNow()
	}
	ingredientRepository := ingredientrepo.New(db)
	ctx := context.Background()
	return ingredientRepository, clock, ctx
}

func TestAddIngredient(t *testing.T) {

	t.Run("should add ID if no error adding ingrediente", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, clock.Now())
		require.NoError(t, ingredientRepository.Add(ctx, ingredient))
		assert.Equal(t, int64(1), ingredient.ID)
	})

	t.Run("should not add a valid ID if error", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, clock.Now())
		require.NoError(t, ingredientRepository.Add(ctx, ingredient))
		assert.Equal(t, int64(1), ingredient.ID)
	})
}

func TestFindIngredient(t *testing.T) {

	t.Run("should find correct ingredient if existent", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)
		now := clock.Now()
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, now)
		require.NoError(t, ingredientRepository.Add(ctx, ingredient))
		ingredient2, _ := model.NewIngredient("ing2", model.Gram, 1234.0, now)
		require.NoError(t, ingredientRepository.Add(ctx, ingredient2))
		ingr1Get, err := ingredientRepository.Find(ctx, ingredient.ID)
		require.NoError(t, err)
		assert.Equal(t, ingredient, &ingr1Get)
	})

	t.Run("should return error when finding an unexistent ingredient", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, clock.Now())
		_, err := ingredientRepository.Find(ctx, ingredient.ID)
		require.Error(t, err)
		assert.Equal(t, err, errs.ErrNotFound)
		assert.Equal(t, int64(-1), ingredient.ID)
	})

	t.Run("should return internal error if returned", func(t *testing.T) {
		db := new(mocks.DatabaseMock)
		ingredientRepository := ingredientrepo.New(db)
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, clock.New().Now())
		_, err := ingredientRepository.Find(context.Background(), ingredient.ID)
		require.Error(t, err)
		assert.Equal(t, err, mocks.ErrDBInternal)
	})
}

func TestUpdateIngredient(t *testing.T) {

	t.Run("should find correct ingredient if existent", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)
		now := clock.Now()
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, now)
		require.NoError(t, ingredientRepository.Add(ctx, ingredient))
		ingredientRepository.Update(ctx, ingredient.ID, func(ingredient *model.Ingredient) error {
			ingredient.Price = 2.0
			ingredient.Name = "new name"
			return nil
		})
		ingr1Get, err := ingredientRepository.Find(ctx, ingredient.ID)
		require.NoError(t, err)
		assert.Equal(t, "new name", ingr1Get.Name)
		assert.Equal(t, 2.0, ingr1Get.Price)
	})

	t.Run("should not update ingredient if ingredient not found", func(t *testing.T) {
		ingredientRepository, _, ctx := setupTest(t)
		err := ingredientRepository.Update(ctx, 1.0, func(ingredient *model.Ingredient) error {
			ingredient.Price = 2.0
			ingredient.Name = "new name"
			return nil
		})
		require.Error(t, err, errs.ErrNotFound)
	})

	t.Run("should not update ingredient if error found", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, clock.Now())
		require.NoError(t, ingredientRepository.Add(ctx, ingredient))
		require.Error(t, ingredientrepo.New(new(mocks.DatabaseMock)).Update(ctx, ingredient.ID, func(ingredient *model.Ingredient) error {
			ingredient.Price = 2.0
			ingredient.Name = "new name"
			return nil
		}))

		ingr1Get, err := ingredientRepository.Find(context.Background(), ingredient.ID)
		require.NoError(t, err)
		assert.Equal(t, ingredient.Price, ingr1Get.Price)
		assert.Equal(t, ingredient.Name, ingr1Get.Name)
	})
}

func TestFindAllIngredients(t *testing.T) {

	t.Run("should get list of existent ingredients", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)
		now := clock.Now()
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, now)
		err := ingredientRepository.Add(ctx, ingredient)
		require.NoError(t, err)
		ingredient2, _ := model.NewIngredient("ing2", model.Gram, 1234.0, now)
		err = ingredientRepository.Add(ctx, ingredient2)
		require.NoError(t, err)

		ingredients, err := ingredientRepository.FindAll(ctx)
		require.NoError(t, err)

		assert.Equal(t, ingredient, &ingredients[0])
		assert.Equal(t, ingredient2, &ingredients[1])
	})
}

func TestIncreaseStockAndUpdatePrice(t *testing.T) {

	t.Run("should increase stock if existent", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)
		now := clock.Now()
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, now)
		require.NoError(t, ingredientRepository.Add(ctx, ingredient))

		ingredientRepository.IncreaseStockAndUpdatePrice(ctx, ingredient.ID, 2, 10.0, clock.Now())
		modifiedTime := clock.Now()
		ingredientRepository.IncreaseStockAndUpdatePrice(ctx, ingredient.ID, 4, 12.0, modifiedTime)

		ingr1Get, err := ingredientRepository.Find(ctx, ingredient.ID)
		require.NoError(t, err)

		assert.Equal(t, 6, ingr1Get.UnitsInStock)
		assert.Equal(t, 12.0, ingr1Get.Price)
		assert.Equal(t, modifiedTime, ingr1Get.LastModified)
	})

	t.Run("should return error not found if unexistent", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)

		err := ingredientRepository.IncreaseStockAndUpdatePrice(ctx, 1.0, 2, 10.0, clock.Now())

		assert.Equal(t, err, errs.ErrNotFound)
	})
}

func TestDecreaseStock(t *testing.T) {

	t.Run("should increase stock if existent", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)
		now := clock.Now()
		ingredient, _ := model.NewIngredient("ing1", model.Gram, 1.0, now)
		require.NoError(t, ingredientRepository.Add(ctx, ingredient))

		modifiedTime := clock.Now()
		ingredientRepository.DecreaseStock(ctx, ingredient.ID, 4, modifiedTime)

		ingr1Get, err := ingredientRepository.Find(ctx, ingredient.ID)
		require.NoError(t, err)

		assert.Equal(t, ingredient.UnitsInStock-4, ingr1Get.UnitsInStock)
		assert.Equal(t, modifiedTime, ingr1Get.LastModified)
	})

	t.Run("should return error not found if unexistent", func(t *testing.T) {
		ingredientRepository, clock, ctx := setupTest(t)

		err := ingredientRepository.DecreaseStock(ctx, 1.0, 2, clock.Now())

		assert.Equal(t, err, errs.ErrNotFound)
	})
}
