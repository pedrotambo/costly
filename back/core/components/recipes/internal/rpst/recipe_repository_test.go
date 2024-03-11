package rpst_test

import (
	"context"
	"costly/core/components/ingredients"
	"costly/core/components/recipes/internal/rpst"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createDBWithIngredients(logger logger.Logger, clock clock.Clock) database.Database {
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientComponent := ingredients.New(db, clock, logger)
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
		ingredientComponent.Create(ctx, ingredient)
	}

	return db
}

func TestGetRecipe(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get correct recipe if existent", func(t *testing.T) {
		db := createDBWithIngredients(logger, clock)

		repo := rpst.New(db, logger)
		ctx := context.Background()
		now := clock.Now()
		recipe1 := model.Recipe{
			ID:           -1,
			Name:         "aName",
			Ingredients:  []model.RecipeIngredient{},
			CreatedAt:    now,
			LastModified: now,
		}
		err := repo.Add(ctx, &recipe1)
		require.NoError(t, err)
		recipe2 := model.Recipe{
			ID:           -1,
			Name:         "anotherName",
			Ingredients:  []model.RecipeIngredient{},
			CreatedAt:    now,
			LastModified: now,
		}
		err = repo.Add(ctx, &recipe2)
		require.NoError(t, err)

		recipe1Get, err := repo.Find(ctx, recipe1.ID)
		require.NoError(t, err)

		assert.Equal(t, recipe1, recipe1Get)
	})
}

func TestGetRecipes(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get correct recipes if existent", func(t *testing.T) {
		db := createDBWithIngredients(logger, clock)

		repo := rpst.New(db, logger)
		ctx := context.Background()
		now := clock.Now()
		recipe1 := model.Recipe{
			ID:           -1,
			Name:         "aName",
			Ingredients:  []model.RecipeIngredient{},
			CreatedAt:    now,
			LastModified: now,
		}
		err := repo.Add(ctx, &recipe1)
		require.NoError(t, err)
		recipe2 := model.Recipe{
			ID:           -1,
			Name:         "anotherName",
			Ingredients:  []model.RecipeIngredient{},
			CreatedAt:    now,
			LastModified: now,
		}
		err = repo.Add(ctx, &recipe2)
		require.NoError(t, err)

		recipes, err := repo.FindAll(ctx)
		require.NoError(t, err)

		assert.Equal(t, recipe1, recipes[0])
		assert.Equal(t, recipe2, recipes[1])
	})
}
