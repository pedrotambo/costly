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

func TestHandleGetRecipe(t *testing.T) {
	clock := new(mocks.ClockMock)
	now := time.UnixMilli(12345).UTC()
	clock.On("Now").Return(now)

	testCases := []struct {
		name        string
		recipeIDstr string
		expected    string
		statusCode  int
	}{
		{
			name:        "should create recipe if id is valid",
			recipeIDstr: "1",
			expected: `{
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
			}`,
			statusCode: http.StatusOK,
		},
		{
			name:        "should get error if unexistent recipe",
			recipeIDstr: "123",
			expected:    "",
			statusCode:  http.StatusNotFound,
		},
		{
			name:        "should get error if bad request id",
			recipeIDstr: "badID",
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"id is invalid"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/recipes/"+tc.recipeIDstr, nil)
			require.NoError(t, err)
			rr := makeRequest(t, clock, func(useCases *usecases.UseCases) error {
				_, err := useCases.Ingredients.Create(context.Background(), ingredients.CreateIngredientOptions{
					Name:  "ingr1",
					Price: 1.50,
					Unit:  model.Gram,
				})
				require.NoError(t, err)
				_, err = useCases.Ingredients.Create(context.Background(), ingredients.CreateIngredientOptions{
					Name:  "ingr2",
					Price: 2.50,
					Unit:  model.Gram,
				})
				require.NoError(t, err)
				_, err = useCases.Recipes.Create(context.Background(), recipes.CreateRecipeOptions{
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
				})
				require.NoError(t, err)
				return nil
			}, req)
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
