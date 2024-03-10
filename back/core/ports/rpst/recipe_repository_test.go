package rpst_test

import (
	"context"
	"costly/core/domain"
	clck "costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var now = clck.New().Now()

var meat = &domain.Ingredient{
	ID:           -1,
	Name:         "meat",
	Unit:         domain.Gram,
	Price:        1.0,
	UnitsInStock: 0,
	CreatedAt:    now,
	LastModified: now,
}

var salt = &domain.Ingredient{
	ID:           -1,
	Name:         "salt",
	Price:        10.0,
	Unit:         domain.Gram,
	UnitsInStock: 0,
	CreatedAt:    now,
	LastModified: now,
}

var pepper = &domain.Ingredient{
	ID:           -1,
	Name:         "pepper",
	Price:        13.0,
	Unit:         domain.Gram,
	UnitsInStock: 0,
	CreatedAt:    now,
	LastModified: now,
}

func createDBWithIngredients(logger logger.Logger, clock clck.Clock) database.Database {
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)
	ctx := context.Background()

	var ingredients = []*domain.Ingredient{meat, salt, pepper}
	for _, ingredient := range ingredients {
		ingredientRepository.SaveIngredient(ctx, ingredient)
	}

	return db
}

func TestGetRecipe(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clck.New()

	t.Run("should get correct recipe if existent", func(t *testing.T) {
		db := createDBWithIngredients(logger, clock)

		repo := rpst.NewRecipeRepository(db, clock, logger)
		ctx := context.Background()
		now := clock.Now()
		recipe1 := domain.Recipe{
			ID:           -1,
			Name:         "aName",
			Ingredients:  []domain.RecipeIngredient{},
			CreatedAt:    now,
			LastModified: now,
		}
		err := repo.SaveRecipe(ctx, &recipe1)
		require.NoError(t, err)
		recipe2 := domain.Recipe{
			ID:           -1,
			Name:         "anotherName",
			Ingredients:  []domain.RecipeIngredient{},
			CreatedAt:    now,
			LastModified: now,
		}
		err = repo.SaveRecipe(ctx, &recipe2)
		require.NoError(t, err)

		recipe1Get, err := repo.GetRecipe(ctx, recipe1.ID)
		require.NoError(t, err)

		assert.Equal(t, recipe1, recipe1Get)
	})
}

func TestGetRecipes(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clck.New()

	t.Run("should get correct recipes if existent", func(t *testing.T) {
		db := createDBWithIngredients(logger, clock)

		repo := rpst.New(db, clock, logger)
		ctx := context.Background()
		now := clock.Now()
		recipe1 := domain.Recipe{
			ID:           -1,
			Name:         "aName",
			Ingredients:  []domain.RecipeIngredient{},
			CreatedAt:    now,
			LastModified: now,
		}
		err := repo.SaveRecipe(ctx, &recipe1)
		require.NoError(t, err)
		recipe2 := domain.Recipe{
			ID:           -1,
			Name:         "anotherName",
			Ingredients:  []domain.RecipeIngredient{},
			CreatedAt:    now,
			LastModified: now,
		}
		err = repo.SaveRecipe(ctx, &recipe2)
		require.NoError(t, err)

		recipes, err := repo.GetRecipes(ctx)
		require.NoError(t, err)

		assert.Equal(t, recipe1, recipes[0])
		assert.Equal(t, recipe2, recipes[1])
	})
}
