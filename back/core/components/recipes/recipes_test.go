package recipes_test

import (
	"context"
	"costly/core/components/clock"
	"costly/core/components/database"
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/components/recipes"
	"costly/core/errs"
	"costly/core/mocks"
	"costly/core/model"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var meat = ingredients.CreateIngredientOptions{
	Name:  "meat",
	Price: 1.0,
	Unit:  model.Gram,
}

var salt = ingredients.CreateIngredientOptions{
	Name:  "salt",
	Price: 10.0,
	Unit:  model.Gram,
}

var pepper = ingredients.CreateIngredientOptions{
	Name:  "pepper",
	Price: 13.0,
	Unit:  model.Gram,
}

var ingredientOptsByName = map[string]ingredients.CreateIngredientOptions{meat.Name: meat, salt.Name: salt, pepper.Name: pepper}

func createDBWithIngredients(logger logger.Logger, clock clock.Clock) (database.Database, []model.Ingredient, ingredients.IngredientComponent) {
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientComponent := ingredients.New(db, clock, logger)
	ctx := context.Background()
	var createdIngredients = []model.Ingredient{}
	for _, ingredient := range []ingredients.CreateIngredientOptions{meat, salt, pepper} {
		ing, _ := ingredientComponent.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  ingredient.Name,
			Price: ingredient.Price,
			Unit:  ingredient.Unit,
		})
		createdIngredients = append(createdIngredients, *ing)
	}
	return db, createdIngredients, ingredientComponent
}

func TestCreateRecipe(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()
	now := clock.Now()

	t.Run("should create a recipe if non existent", func(t *testing.T) {
		logger.Debug(now.String())
		clockMock := new(mocks.ClockMock)
		now := time.UnixMilli(12345).UTC()
		clockMock.On("Now").Return(now)
		db, ingredients, ingredientComponent := createDBWithIngredients(logger, clockMock)

		recipeComponent := recipes.New(db, clockMock, logger, ingredientComponent)

		recipe, err := recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name: "recipeName",
			Ingredients: []recipes.RecipeIngredientOptions{
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
		db, ingredients, ingredientComponent := createDBWithIngredients(logger, clock)
		recipeComponent := recipes.New(db, clock, logger, ingredientComponent)
		existentRecipeName := "name"

		recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name: existentRecipeName,
			Ingredients: []recipes.RecipeIngredientOptions{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				},
			},
		})

		_, err := recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name: existentRecipeName,
			Ingredients: []recipes.RecipeIngredientOptions{
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
		db, ingredients, ingredientComponent := createDBWithIngredients(logger, clock)

		recipeComponent := recipes.New(db, clock, logger, ingredientComponent)
		existentRecipeName := "name"
		var unexistentIngredientID int64
		for _, i := range ingredients {
			unexistentIngredientID += i.ID
		}
		_, err := recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name: existentRecipeName,
			Ingredients: []recipes.RecipeIngredientOptions{
				{
					ID:    unexistentIngredientID,
					Units: 500,
				},
			},
		})
		assert.EqualError(t, err, "failed to create recipe: unexistent ingredient with ID 6")
	})

	t.Run("should assign different IDs to different recipes", func(t *testing.T) {
		db, _, ingredientComponent := createDBWithIngredients(logger, clock)

		recipeComponent := recipes.New(db, clock, logger, ingredientComponent)
		recipe1, err := recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name:        "recipe1",
			Ingredients: []recipes.RecipeIngredientOptions{},
		})
		require.NoError(t, err)
		recipe2, err := recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name:        "recipe2",
			Ingredients: []recipes.RecipeIngredientOptions{},
		})
		require.NoError(t, err)

		assert.NotEqual(t, recipe1.ID, recipe2.ID)
	})
}

func TestGetRecipe(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get correct recipe if existent", func(t *testing.T) {
		db, _, ingredientComponent := createDBWithIngredients(logger, clock)

		recipeComponent := recipes.New(db, clock, logger, ingredientComponent)
		ctx := context.Background()
		recipe1, err := recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name:        "recipe1",
			Ingredients: []recipes.RecipeIngredientOptions{},
		})
		require.NoError(t, err)
		_, err = recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name:        "recipe2",
			Ingredients: []recipes.RecipeIngredientOptions{},
		})
		require.NoError(t, err)

		recipe1Get, err := recipeComponent.Find(ctx, recipe1.ID)
		require.NoError(t, err)

		assert.Equal(t, recipe1, &recipe1Get)
	})

	t.Run("should return error when requesting an inexistent ingredient", func(t *testing.T) {
		db, _, ingredientComponent := createDBWithIngredients(logger, clock)
		recipeComponent := recipes.New(db, clock, logger, ingredientComponent)

		_, err := recipeComponent.Find(context.Background(), 123)

		require.Error(t, err)
		assert.Equal(t, err, errs.ErrNotFound)
	})
}

func TestGetRecipes(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get correct recipes if existent", func(t *testing.T) {
		db, _, ingredientComponent := createDBWithIngredients(logger, clock)

		recipeComponent := recipes.New(db, clock, logger, ingredientComponent)
		ctx := context.Background()
		recipe1, err := recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name:        "recipe1",
			Ingredients: []recipes.RecipeIngredientOptions{},
		})
		require.NoError(t, err)
		recipe2, err := recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name:        "recipe2",
			Ingredients: []recipes.RecipeIngredientOptions{},
		})
		require.NoError(t, err)

		recipes, err := recipeComponent.FindAll(ctx)
		require.NoError(t, err)

		assert.Equal(t, recipe1, &recipes[0])
		assert.Equal(t, recipe2, &recipes[1])
	})
}
