package model_test

import (
	"costly/core/model"
	"costly/core/ports/clock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeCost(t *testing.T) {

	t.Run("cost of a recipe is the sum of its ingredients and units of them", func(t *testing.T) {
		clock := clock.New()
		now := clock.Now()
		recipe := model.RecipeView{
			ID:   1,
			Name: "aName",
			Ingredients: []model.RecipeIngredientView{
				{
					ID:    1,
					Name:  "meat",
					Units: 500,
					Price: 1,
				},
				{
					ID:    2,
					Name:  "salt",
					Units: 5,
					Price: 10.0,
				},
			},
			CreatedAt:    now,
			LastModified: now,
		}

		assert.Equal(t, recipe.Cost(), 500.0+5*10)
	})
}
