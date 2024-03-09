package handlers_test

import (
	"bytes"
	"context"
	"costly/api/handlers"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func runCreateRecipeHandler(t *testing.T, clock clock.Clock, reqBody io.Reader) *httptest.ResponseRecorder {
	logger, _ := logger.New("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientRepository := rpst.NewIngredientRepository(db, clock, logger)
	ingredientRepository.CreateIngredient(context.Background(), rpst.CreateIngredientOptions{
		Name:  "ingr1",
		Price: 1.50,
		Unit:  domain.Gram,
	})
	ingredientRepository.CreateIngredient(context.Background(), rpst.CreateIngredientOptions{
		Name:  "ingr2",
		Price: 2.50,
		Unit:  domain.Gram,
	})

	repo := rpst.NewRecipeRepository(db, clock, logger)
	handler := handlers.CreateRecipeHandler(repo)

	req, err := http.NewRequest("POST", "/recipes", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/recipes", handler)
	mux.ServeHTTP(rr, req)

	return rr
}

func TestHandleCreateRecipe(t *testing.T) {
	clock := new(clockMock)
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
						"ingredient": {
							"id": 1,
							"name": "ingr1",
							"unit": "gr",
							"price": 1.50,
							"units_in_stock":0,
							"created_at": "1970-01-01T00:00:12.345Z",
							"last_modified": "1970-01-01T00:00:12.345Z"
						},
						"units": 5
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
						"units": 500
					}
				],
				"created_at": "1970-01-01T00:00:12.345Z",
				"last_modified": "1970-01-01T00:00:12.345Z",
				"cost": 1257.5
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
			rr := runCreateRecipeHandler(t, clock, bytes.NewBufferString(tc.payload))
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
