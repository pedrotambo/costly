package recipes_test

import (
	"context"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/logger"
	"costly/core/usecases/recipes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindAll(t *testing.T) {
	logger, _ := logger.New("debug")
	clock := clock.New()

	t.Run("should get correct recipes if existent", func(t *testing.T) {
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
		recipe2, err := recipeComponent.Create(context.Background(), recipes.CreateRecipeOptions{
			Name: "recipe2",
			Ingredients: []model.RecipeIngredient{
				{
					ID:    ingredients[0].ID,
					Units: 500,
				},
			},
		})
		require.NoError(t, err)

		recipes, err := recipeComponent.FindAll(ctx)
		require.NoError(t, err)
		assert.Equal(t, recipe1.Name, recipes[0].Name)
		assert.Equal(t, recipe1.ID, recipes[0].ID)
		assert.Equal(t, recipe1.CreatedAt, recipes[0].CreatedAt)
		assert.Equal(t, recipe1.LastModified, recipes[0].LastModified)

		assert.Equal(t, recipe2.Name, recipes[1].Name)
		assert.Equal(t, recipe2.ID, recipes[1].ID)
		assert.Equal(t, recipe2.CreatedAt, recipes[1].CreatedAt)
		assert.Equal(t, recipe2.LastModified, recipes[1].LastModified)
	})
}
