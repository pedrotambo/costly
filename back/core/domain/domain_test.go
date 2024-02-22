package domain_test

import (
	"costly/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEntityEquality(t *testing.T) {

	t.Run("cost of a recipe is the sum of its ingredients and units of them", func(t *testing.T) {

		// type Ingredient struct {
		// 	ID           int64     `json:"id"`
		// 	Name         string    `json:"name"`
		// 	Unit         Unit      `json:"unit"`
		// 	Price        float64   `json:"price"`
		// 	CreatedAt    time.Time `json:"created_at"`
		// 	LastModified time.Time `json:"last_modified"`
		// }
		now := time.Now()
		i1 := domain.Ingredient{
			1,
			"perro",
			domain.Gram,
			10.0,
			now,
			now,
		}
		i2 := domain.Ingredient{
			1,
			"perro",
			domain.Gram,
			10.0,
			now,
			now,
		}

		assert.Equal(t, true, i1 == i2)

		// type RecipeIngredient struct {
		// 	Ingredient Ingredient `json:"ingredient"`
		// 	Units      int        `json:"units"`
		// }

		// type Recipe struct {
		// 	ID           int64              `json:"id"`
		// 	Name         string             `json:"name"`
		// 	Ingredients  []RecipeIngredient `json:"ingredients"`
		// 	CreatedAt    time.Time          `json:"created_at"`
		// 	LastModified time.Time          `json:"last_modified"`
		// }
		r1 := domain.Recipe{
			1,
			"perro",
			[]domain.RecipeIngredient{},
			now,
			now,
		}
		r2 := domain.Recipe{
			1,
			"perro",
			[]domain.RecipeIngredient{},
			now,
			now,
		}

		// assert.Equal(t, i1, i2)
		assert.Equal(t, true, r1 == r2)
	})
}
