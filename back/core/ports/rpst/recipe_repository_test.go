package rpst_test

import (
	"context"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"errors"
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var meat = rpst.CreateIngredientOptions{
	Name:  "meat",
	Price: 1.0,
	Unit:  domain.Gram,
}

var salt = rpst.CreateIngredientOptions{
	Name:  "salt",
	Price: 10.0,
	Unit:  domain.Gram,
}

var pepper = rpst.CreateIngredientOptions{
	Name:  "pepper",
	Price: 13.0,
	Unit:  domain.Gram,
}

var ingredientOptsByName = map[string]rpst.CreateIngredientOptions{meat.Name: meat, salt.Name: salt, pepper.Name: pepper}

func createDBWithIngredients(logger logger.Logger, clock clock.Clock) (database.Database, []domain.Ingredient) {
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)
	ctx := context.Background()
	var ingredients = []domain.Ingredient{}
	for _, ingredient := range []rpst.CreateIngredientOptions{meat, salt, pepper} {
		ing, _ := ingredientRepository.CreateIngredient(ctx, rpst.CreateIngredientOptions{
			Name:  ingredient.Name,
			Price: ingredient.Price,
			Unit:  ingredient.Unit,
		})
		ingredients = append(ingredients, ing)
	}
	return db, ingredients
}

func TestCreateRecipe(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()
	now := clock.Now()

	t.Run("should create a recipe if non existent", func(t *testing.T) {
		logger.Debug(now.String())
		clockMock := new(clockMock)
		now := time.UnixMilli(12345).UTC()
		clockMock.On("Now").Return(now)
		db, ingredients := createDBWithIngredients(logger, clockMock)
		fmt.Println(ingredients)

		recipeRepository := rpst.NewRecipeRepository(db, clockMock, logger)

		recipe, err := recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name: "recipeName",
			Ingredients: []rpst.RecipeIngredientOptions{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				}, {
					ID:    ingredients[2].ID,
					Units: 5,
				},
			},
		})

		if err != nil {
			logger.Error(err, "error")
			t.Fail()
		}

		assert.Equal(t, int64(1), recipe.ID)
		assert.Equal(t, "recipeName", recipe.Name)
		assert.Len(t, recipe.Ingredients, 2)
		for _, recipeIngredient := range recipe.Ingredients {
			ingredientOpts := ingredientOptsByName[recipeIngredient.Ingredient.Name]
			assert.Equal(t, ingredientOpts.Unit, recipeIngredient.Ingredient.Unit)
			assert.Equal(t, ingredientOpts.Price, recipeIngredient.Ingredient.Price)
			if recipeIngredient.Ingredient.Name == "meat" {
				assert.Equal(t, 500, recipeIngredient.Units)
			} else {
				assert.Equal(t, 5, recipeIngredient.Units)
			}

		}
		assert.Equal(t, recipe.CreatedAt, now)
		assert.Equal(t, recipe.LastModified, now)
	})

	t.Run("should fail to create a recipe if existent", func(t *testing.T) {
		db, ingredients := createDBWithIngredients(logger, clock)

		recipeRepository := rpst.NewRecipeRepository(db, clock, logger)
		existentRecipeName := "name"

		recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name: existentRecipeName,
			Ingredients: []rpst.RecipeIngredientOptions{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				},
			},
		})

		_, err := recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name: existentRecipeName,
			Ingredients: []rpst.RecipeIngredientOptions{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				},
			},
		})
		require.Error(t, err)
		assert.EqualError(t, err, "failed to create recipe: UNIQUE constraint failed: recipe.name")
	})

	t.Run("should return an error when creating a recipe with unexistent ingredient", func(t *testing.T) {
		db, ingredients := createDBWithIngredients(logger, clock)

		recipeRepository := rpst.NewRecipeRepository(db, clock, logger)
		existentRecipeName := "name"
		var unexistentIngredientID int64
		for _, i := range ingredients {
			unexistentIngredientID += i.ID
		}
		_, err := recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name: existentRecipeName,
			Ingredients: []rpst.RecipeIngredientOptions{
				{
					ID:    unexistentIngredientID,
					Units: 500,
				},
			},
		})
		assert.Equal(t, rpst.ErrBadOpts, errors.Unwrap(err))
		assert.EqualError(t, err, "failed to create recipe: bad create entity options")
	})

	t.Run("should assign different IDs to different recipes", func(t *testing.T) {
		db, _ := createDBWithIngredients(logger, clock)

		recipeRepository := rpst.NewRecipeRepository(db, clock, logger)
		recipe1, err := recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name:        "recipe1",
			Ingredients: []rpst.RecipeIngredientOptions{},
		})
		require.NoError(t, err)
		recipe2, err := recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name:        "recipe2",
			Ingredients: []rpst.RecipeIngredientOptions{},
		})
		require.NoError(t, err)

		assert.NotEqual(t, recipe1.ID, recipe2.ID)
	})
}

func TestGetRecipe(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get correct recipe if existent", func(t *testing.T) {
		db, _ := createDBWithIngredients(logger, clock)

		recipeRepository := rpst.NewRecipeRepository(db, clock, logger)
		ctx := context.Background()
		recipe1, err := recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name:        "recipe1",
			Ingredients: []rpst.RecipeIngredientOptions{},
		})
		require.NoError(t, err)
		_, err = recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name:        "recipe2",
			Ingredients: []rpst.RecipeIngredientOptions{},
		})
		require.NoError(t, err)

		recipe1Get, err := recipeRepository.GetRecipe(ctx, recipe1.ID)
		require.NoError(t, err)

		assert.Equal(t, recipe1, recipe1Get)
	})

	t.Run("should return error when requesting an inexistent ingredient", func(t *testing.T) {
		db, _ := createDBWithIngredients(logger, clock)
		recipeRepository := rpst.NewRecipeRepository(db, clock, logger)

		_, err := recipeRepository.GetRecipe(context.Background(), 123)

		require.Error(t, err)
		assert.Equal(t, err, rpst.ErrNotFound)
	})
}

func TestGetRecipes(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get correct recipes if existent", func(t *testing.T) {
		db, _ := createDBWithIngredients(logger, clock)

		recipeRepository := rpst.NewRecipeRepository(db, clock, logger)
		ctx := context.Background()
		recipe1, err := recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name:        "recipe1",
			Ingredients: []rpst.RecipeIngredientOptions{},
		})
		require.NoError(t, err)
		recipe2, err := recipeRepository.CreateRecipe(context.Background(), rpst.CreateRecipeOptions{
			Name:        "recipe2",
			Ingredients: []rpst.RecipeIngredientOptions{},
		})
		require.NoError(t, err)

		recipes, err := recipeRepository.GetRecipes(ctx)
		require.NoError(t, err)

		assert.Equal(t, recipe1, recipes[0])
		assert.Equal(t, recipe2, recipes[1])
	})
}
