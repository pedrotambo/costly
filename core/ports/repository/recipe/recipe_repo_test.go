package reciperepo_test

import (
	"context"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	reciperepo "costly/core/ports/repository/recipe"
	"costly/core/usecases/ingredients"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createDBWithIngredients(t *testing.T, logger logger.Logger, clock clock.Clock) database.Database {
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientUseCases := ingredients.New(db, clock)
	ctx := context.Background()
	var ingredientOpts = []ingredients.CreateIngredientOptions{
		{
			Name:  "meat",
			Price: 1.0,
			Unit:  model.Gram,
		},
		{
			Name:  "salt",
			Price: 10.0,
			Unit:  model.Gram,
		},
		{
			Name:  "pepper",
			Price: 13.0,
			Unit:  model.Gram,
		},
	}
	for _, ingredient := range ingredientOpts {
		if _, err := ingredientUseCases.Create(ctx, ingredient); err != nil {
			t.FailNow()
		}
	}
	return db
}

func TestFind(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get correct recipe if existent", func(t *testing.T) {
		db := createDBWithIngredients(t, logger, clock)

		repo := reciperepo.New(db)
		ctx := context.Background()
		now := clock.Now()

		recipe1, err := model.NewRecipe("aName", []model.RecipeIngredient{{ID: 1, Units: 1}}, now)
		require.NoError(t, err)
		require.NoError(t, repo.Add(ctx, recipe1))
		recipe2, err := model.NewRecipe("anotherName", []model.RecipeIngredient{{ID: 1, Units: 1}}, now)
		require.NoError(t, err)
		require.NoError(t, repo.Add(ctx, recipe2))
		require.NoError(t, err)
		recipe1Get, err := repo.Find(ctx, recipe1.ID)
		require.NoError(t, err)
		assert.Equal(t, *recipe1, recipe1Get)
	})
}
