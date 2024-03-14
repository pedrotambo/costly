package ingredients_test

import (
	"costly/core/model"
	"costly/core/usecases/ingredients"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {

	t.Run("should edit ingredient correctly", func(t *testing.T) {

		ingredientComponent, ctx := setupTest(t)
		ing1, err := ingredientComponent.Create(ctx, ingredients.CreateIngredientOptions{
			Name:  "ing1",
			Price: 10.0,
			Unit:  model.Gram,
		})
		require.NoError(t, err)

		newIngredientOpts := ingredients.CreateIngredientOptions{
			Name:  "modifiedIngr1",
			Price: ing1.Price + 10.0,
			Unit:  model.Gram,
		}
		err = ingredientComponent.Update(ctx, ing1.ID, newIngredientOpts)
		require.NoError(t, err)

		modifiedIngredient, err := ingredientComponent.Find(ctx, ing1.ID)
		require.NoError(t, err)

		assert.Equal(t, modifiedIngredient.Name, newIngredientOpts.Name)
		assert.Equal(t, modifiedIngredient.Price, newIngredientOpts.Price)
		assert.Equal(t, modifiedIngredient.Unit, newIngredientOpts.Unit)
	})
}
