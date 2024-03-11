package handlers_test

import (
	"context"
	"costly/api/handlers"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func runGetRecipesHandler(t *testing.T, clock clock.Clock, recipeOpts []usecases.CreateRecipeOptions) *httptest.ResponseRecorder {
	logger, _ := logger.New("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	useCases := usecases.New(&usecases.Ports{
		Logger:     logger,
		Repository: rpst.New(db, clock, logger),
		Clock:      clock,
	})
	useCases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
		Name:  "ingr1",
		Price: 1.50,
		Unit:  model.Gram,
	})
	useCases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
		Name:  "ingr2",
		Price: 2.50,
		Unit:  model.Gram,
	})

	for _, opts := range recipeOpts {
		useCases.CreateRecipe(context.Background(), opts)
	}

	handler := handlers.GetRecipesHandler(useCases)

	req, err := http.NewRequest("GET", "/ingredients", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/ingredients", handler)
	mux.ServeHTTP(rr, req)

	return rr
}

func TestHandleGetRecipes(t *testing.T) {
	clock := new(clockMock)
	now := time.UnixMilli(12345).UTC()
	clock.On("Now").Return(now)

	testCases := []struct {
		name       string
		recipes    []usecases.CreateRecipeOptions
		expected   string
		statusCode int
	}{
		{
			name: "should get recipes",
			recipes: []usecases.CreateRecipeOptions{
				{
					Name: "recipe1",
					Ingredients: []usecases.RecipeIngredientOptions{
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
					Ingredients: []usecases.RecipeIngredientOptions{
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
							"ingredient": {
								"id": 1,
								"name": "ingr1",
								"unit": "gr",
								"price": 1.50,
								"units_in_stock": 0,
								"created_at": "1970-01-01T00:00:12.345Z",
								"last_modified": "1970-01-01T00:00:12.345Z"
							},
							"units": 1
						},
						{
							"ingredient": {
								"id": 2,
								"name": "ingr2",
								"unit": "gr",
								"price": 2.50,
								"units_in_stock": 0,
								"created_at": "1970-01-01T00:00:12.345Z",
								"last_modified": "1970-01-01T00:00:12.345Z"
							},
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
							"ingredient": {
								"id": 2,
								"name": "ingr2",
								"unit": "gr",
								"price": 2.50,
								"units_in_stock": 0,
								"created_at": "1970-01-01T00:00:12.345Z",
								"last_modified": "1970-01-01T00:00:12.345Z"
							},
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
			recipes:    []usecases.CreateRecipeOptions{},
			expected:   `[]`,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := runGetRecipesHandler(t, clock, tc.recipes)
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
