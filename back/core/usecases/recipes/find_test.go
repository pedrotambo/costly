package recipes_test

import (
	"costly/core/errs"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/logger"
	"costly/core/usecases/recipes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get correct recipe if existent", func(t *testing.T) {
		ingredients, recipeComponent, ctx := setupTest(logger, clock)
		recipe1, err := recipeComponent.Create(ctx, recipes.CreateRecipeOptions{
			Name: "recipe1",
			Ingredients: []model.RecipeIngredient{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				},
			},
		})
		require.NoError(t, err)
		_, err = recipeComponent.Create(ctx, recipes.CreateRecipeOptions{
			Name: "recipe2",
			Ingredients: []model.RecipeIngredient{
				{
					ID:    ingredients[1].ID,
					Units: 25,
				},
			},
		})
		require.NoError(t, err)

		recipe1Get, err := recipeComponent.Find(ctx, recipe1.ID)
		require.NoError(t, err)
		assert.Equal(t, recipe1.Name, recipe1Get.Name)
		assert.Equal(t, recipe1.ID, recipe1Get.ID)
		assert.Equal(t, recipe1.CreatedAt, recipe1Get.CreatedAt)
		assert.Equal(t, recipe1.LastModified, recipe1Get.LastModified)
		assert.Equal(t, ingredients[0].ID, recipe1Get.Ingredients[0].ID)
		assert.Equal(t, 500, recipe1Get.Ingredients[0].Units)

	})

	t.Run("should return error when requesting an inexistent ingredient", func(t *testing.T) {
		_, recipeComponent, ctx := setupTest(logger, clock)
		_, err := recipeComponent.Find(ctx, 123)

		require.Error(t, err)
		assert.Equal(t, err, errs.ErrNotFound)
	})
}
