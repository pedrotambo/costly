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

func TestHandleGetIngredient(t *testing.T) {
	clock := new(mocks.ClockMock)
	now := time.UnixMilli(12345).UTC()
	clock.On("Now").Return(now)

	testCases := []struct {
		name            string
		ingredientIDstr string
		expected        string
		statusCode      int
	}{
		{
			name:            "should create ingredient if payload is valid",
			ingredientIDstr: "1",
			expected: `{
				"id":1,
				"name":"ingredientName",
				"unit":"gr",
				"price":12.43,
				"units_in_stock":0,
				"created_at":"1970-01-01T00:00:12.345Z",
				"last_modified":"1970-01-01T00:00:12.345Z"
			}`,
			statusCode: http.StatusOK,
		},
		{
			name:            "should get error if unexistent ingredient",
			ingredientIDstr: "123",
			expected:        "",
			statusCode:      http.StatusNotFound,
		},
		{
			name:            "should get error if bad request id",
			ingredientIDstr: "badID",
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
			req, err := http.NewRequest("GET", "/ingredients/"+tc.ingredientIDstr, nil)
			require.NoError(t, err)
			rr := makeRequest(t, clock, func(components *comps.Components) error {
				components.Ingredients.Create(context.Background(), ingredients.CreateIngredientOptions{
					Name:  "ingredientName",
					Price: 12.43,
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
