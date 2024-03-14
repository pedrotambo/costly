package handlers_test

import (
	"bytes"
	"context"
	"costly/api"
	"costly/core/mocks"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/usecases"
	"costly/core/usecases/ingredients"
	"costly/core/usecases/recipes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var dummyHandler = api.Middleware(func(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
})

func makeRequest(t *testing.T, clock clock.Clock, prepare func(useCases *usecases.UseCases) error, req *http.Request) *httptest.ResponseRecorder {
	logger, _ := logger.New("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientUseCases := ingredients.New(db, clock)
	recipeUseCases := recipes.New(db, clock, logger, ingredientUseCases)
	useCases := &usecases.UseCases{
		Ingredients: ingredientUseCases,
		Recipes:     recipeUseCases,
	}
	err := prepare(useCases)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	var router = api.NewRouter(useCases, dummyHandler)
	router.ServeHTTP(rr, req)
	return rr

}

func TestHandleCreateRecipe(t *testing.T) {
	clock := new(mocks.ClockMock)
	now := time.UnixMilli(12345).UTC()
	clock.On("Now").Return(now)

	testCases := []struct {
		name       string
		payload    string
		expected   string
		statusCode int
	}{
		{
			name: "should create ingredient if payload is valid",
			payload: `{
				"name": "recipe1",
				"ingredients": [
					{
						"id": 1,
						"units": 5
					},
					{
						"id": 2,
						 "units": 500
					}
				]
			}`,
			expected: `{
				"id": 1,
				"name": "recipe1",
				"ingredients": [
					{
						"id": 1,
						"units": 5
					},
					{
						"id": 2,
						"units": 500
					}
				],
				"created_at": "1970-01-01T00:00:12.345Z",
				"last_modified": "1970-01-01T00:00:12.345Z"
			}`,
			statusCode: http.StatusCreated,
		},
		{
			name: "should return error if name is invalid",
			payload: `{
				"name": "",
				"ingredients": [
					{
						"id": 1,
						"units": 5
					},
					{
						"id": 2,
						 "units": 500
					}
				]
			}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"name is invalid"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "should return error if name is valid and ingredients is not present",
			payload: `{
				"name": "validName"
			}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"recipe must have at least one ingredient"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "should return error if name is valid and ingredients is empty",
			payload: `{
				"name": "validName",
				"ingredients": []
			}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"recipe must have at least one ingredient"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:    "should return error if payload is invalid json",
			payload: "invalid payload",
			expected: `{
				"error": {
					"code":"INVALID_JSON",
					"message":"error unmarshalling request body"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/recipes", bytes.NewBufferString(tc.payload))
			if err != nil {
				t.Fatal(err)
			}
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

				return nil
			}, req)
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
