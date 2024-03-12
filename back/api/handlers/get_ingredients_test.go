package handlers_test

import (
	"context"
	comps "costly/core/components"
	"costly/core/components/ingredients"
	"costly/core/mocks"
	"costly/core/model"
	"net/http"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			req, err := http.NewRequest("GET", "/ingredients", nil)
			require.NoError(t, err)
			rr := makeRequest(t, clock, func(components *comps.Components) error {
				for _, opts := range tc.ingredients {
					components.Ingredients.Create(context.Background(), opts)
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
