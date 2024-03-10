package handlers_test

import (
	"context"
	"costly/api/handlers"
	"costly/core/domain"
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

func runGetRecipeHandler(t *testing.T, clock clock.Clock, recipeIDstr string) *httptest.ResponseRecorder {
	logger, _ := logger.New("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	repo := rpst.New(db, clock, logger)
	allUsecases := usecases.New(repo, clock)
	_, err := allUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
		Name:  "ingr1",
		Price: 1.50,
		Unit:  domain.Gram,
	})

	if err != nil {
		t.Fatal()
	}

	_, err = allUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
		Name:  "ingr2",
		Price: 2.50,
		Unit:  domain.Gram,
	})

	if err != nil {
		t.Fatal()
	}

	allUsecases.CreateRecipe(context.Background(), usecases.CreateRecipeOptions{
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
	})

	handler := handlers.GetRecipeHandler(repo)

	req, err := http.NewRequest("GET", "/recipes/"+recipeIDstr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/recipes/{recipeID}", handler)
	mux.ServeHTTP(rr, req)

	return rr
}

func TestHandleGetRecipe(t *testing.T) {
	clock := new(clockMock)
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
						"ingredient": {
							"id": 1,
							"name": "ingr1",
							"unit": "gr",
							"price": 1.50,
							"units_in_stock":0,
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
							"units_in_stock":0,
							"created_at": "1970-01-01T00:00:12.345Z",
							"last_modified": "1970-01-01T00:00:12.345Z"
						},
						"units": 2
					}
				],
				"created_at": "1970-01-01T00:00:12.345Z",
				"last_modified": "1970-01-01T00:00:12.345Z",
				"cost": 6.5
			}`,
			statusCode: http.StatusOK,
		},
		// {
		// 	name:        "should get error if unexistent recipe",
		// 	recipeIDstr: "123",
		// 	expected:    "",
		// 	statusCode:  http.StatusNotFound,
		// },
		// {
		// 	name:        "should get error if bad request id",
		// 	recipeIDstr: "badID",
		// 	expected: `{
		// 		"error": {
		// 			"code":"INVALID_INPUT",
		// 			"message":"id is invalid"
		// 		}
		// 	}`,
		// 	statusCode: http.StatusBadRequest,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := runGetRecipeHandler(t, clock, tc.recipeIDstr)
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
