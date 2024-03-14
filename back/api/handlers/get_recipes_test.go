package handlers_test

import (
	"context"
	"costly/core/mocks"
	"costly/core/model"
	"costly/core/usecases"
	"costly/core/usecases/ingredients"
	"costly/core/usecases/recipes"
	"net/http"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleGetRecipes(t *testing.T) {
	clock := new(mocks.ClockMock)
	now := time.UnixMilli(12345).UTC()
	clock.On("Now").Return(now)

	testCases := []struct {
		name       string
		recipes    []recipes.CreateRecipeOptions
		expected   string
		statusCode int
	}{
		{
			name: "should get recipes",
			recipes: []recipes.CreateRecipeOptions{
				{
					Name: "recipe1",
					Ingredients: []model.RecipeIngredient{
						{
							ID:    1,
							Units: 1,
						},
						{
							ID:    2,
							Units: 2,
						},
					},
				},
				{
					Name: "recipe2",
					Ingredients: []model.RecipeIngredient{
						{
							ID:    2,
							Units: 3,
						},
					},
				},
			},
			expected: `[
				{
					"id": 1,
					"name": "recipe1",
					"ingredients": [
						{
							"id": 1,
							"name": "ingr1",
							"price": 1.50,
							"units": 1
						},
						{
							"id": 2,
							"name": "ingr2",
							"price": 2.50,
							"units": 2
						}
					],
					"created_at": "1970-01-01T00:00:12.345Z",
					"last_modified": "1970-01-01T00:00:12.345Z",
					"cost": 6.5
				},
				{
					"id": 2,
					"name": "recipe2",
					"ingredients": [
						{
							"id": 2,
							"name": "ingr2",
							"price": 2.50,
							"units": 3
						}
					],
					"created_at": "1970-01-01T00:00:12.345Z",
					"last_modified": "1970-01-01T00:00:12.345Z",
					"cost": 7.5
				}
			]`,
			statusCode: http.StatusOK,
		},
		{
			name:       "should get empty ingredients",
			recipes:    []recipes.CreateRecipeOptions{},
			expected:   `[]`,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/recipes", nil)
			require.NoError(t, err)
			rr := makeRequest(t, clock, func(useCases *usecases.UseCases) error {
				useCases.Ingredients.Create(context.Background(), ingredients.CreateIngredientOptions{
					Name:  "ingr1",
					Price: 1.50,
					Unit:  model.Gram,
				})
				useCases.Ingredients.Create(context.Background(), ingredients.CreateIngredientOptions{
					Name:  "ingr2",
					Price: 2.50,
					Unit:  model.Gram,
				})

				for _, opts := range tc.recipes {
					useCases.Recipes.Create(context.Background(), opts)
				}
				return nil
			}, req)
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
