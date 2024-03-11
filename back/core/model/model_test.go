package model_test

import (
	"costly/core/components/clock"
	"costly/core/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeCost(t *testing.T) {

	t.Run("cost of a recipe is the sum of its ingredients and units of them", func(t *testing.T) {
		clock := clock.New()
		now := clock.Now()
		recipe := model.Recipe{
			ID:   1,
			Name: "aName",
			Ingredients: []model.RecipeIngredient{
				{
					Ingredient: model.Ingredient{
						1,
						"meat",
						model.Gram,
						1.0,
						0,
						now,
						now,
					},
					Units: 500,
				},
				{
					Ingredient: model.Ingredient{
						1,
						"salt",
						model.Gram,
						10.0,
						0,
						now,
						now,
					},
					Units: 5,
				},
			},
			CreatedAt:    now,
			LastModified: now,
		}

		assert.Equal(t, recipe.Cost(), 500.0+5*10)
	})
}
