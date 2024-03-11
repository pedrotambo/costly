package handlers_test

import (
	"context"
	"costly/api/handlers"
	"costly/core/components/ingredients"
	"costly/core/mocks"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func runGetIngredientsHandler(t *testing.T, clock clock.Clock, ingrOpts []ingredients.CreateIngredientOptions) *httptest.ResponseRecorder {
	logger, _ := logger.New("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientComponent := ingredients.New(db, clock, logger)
	for _, opts := range ingrOpts {
		ingredientComponent.Create(context.Background(), opts)
	}

	handler := handlers.GetIngredientsHandler(ingredientComponent)

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

func TestHandleGetIngredients(t *testing.T) {
	clock := new(mocks.ClockMock)
	now := time.UnixMilli(12345).UTC()
	clock.On("Now").Return(now)

	testCases := []struct {
		name        string
		ingredients []ingredients.CreateIngredientOptions
		expected    string
		statusCode  int
	}{
		{
			name: "should get ingredients",
			ingredients: []ingredients.CreateIngredientOptions{
				{
					Name:  "ingr1",
					Price: 1.5,
					Unit:  model.Gram,
				},
				{
					Name:  "ingr2",
					Price: 2.5,
					Unit:  model.Gram,
				},
			},
			expected: `[
				{
					"id": 1,
					"name": "ingr1",
					"unit": "gr",
					"price": 1.5,
					"units_in_stock":0,
					"created_at": "1970-01-01T00:00:12.345Z",
					"last_modified": "1970-01-01T00:00:12.345Z"
				},
				{
					"id": 2,
					"name": "ingr2",
					"unit": "gr",
					"price": 2.5,
					"units_in_stock":0,
					"created_at": "1970-01-01T00:00:12.345Z",
					"last_modified": "1970-01-01T00:00:12.345Z"
				}
			]`,
			statusCode: http.StatusOK,
		},
		{
			name:        "should get empty ingredients",
			ingredients: []ingredients.CreateIngredientOptions{},
			expected:    `[]`,
			statusCode:  http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := runGetIngredientsHandler(t, clock, tc.ingredients)
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
