package domain_test

import (
	"costly/core/domain"
	"costly/core/ports/clock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeCost(t *testing.T) {

	t.Run("cost of a recipe is the sum of its ingredients and units of them", func(t *testing.T) {
		clock := clock.New()
		now := clock.Now()
		recipe := domain.Recipe{
			ID:   1,
			Name: "aName",
			Ingredients: []domain.RecipeIngredient{
				{
					Ingredient: domain.Ingredient{
						1,
						"meat",
						domain.Gram,
						1.0,
						0,
						now,
						now,
					},
					Units: 500,
				},
				{
					Ingredient: domain.Ingredient{
						1,
						"salt",
						domain.Gram,
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
