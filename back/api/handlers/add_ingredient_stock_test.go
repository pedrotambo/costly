package handlers_test

import (
	"bytes"
	"context"
	"costly/core/mocks"
	"costly/core/model"
	"costly/core/usecases"
	"costly/core/usecases/ingredients"
	"net/http"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleAddIngredientStock(t *testing.T) {
	clock := new(mocks.ClockMock)
	now := time.UnixMilli(12345).UTC()
	clock.On("Now").Return(now)

	testCases := []struct {
		name            string
		ingredientIDstr string
		payload         string
		expected        string
		statusCode      int
	}{
		{
			name:            "should add ingredient stock",
			ingredientIDstr: "1",
			payload: `{
				"units": 5,
				"price": 12.5
			}`,
			expected: `{
				"id": 1, 
				"ingredient_id":1, 
				"price":12.5, 
				"units":5,
				"created_at":"1970-01-01T00:00:12.345Z"
			}`,
			statusCode: http.StatusCreated,
		},
		{
			name:            "should get error if updating stock if unexistent ingredient",
			ingredientIDstr: "123",
			payload: `{
				"units": 5,
				"price": 12.5
			}`,
			expected:   "",
			statusCode: http.StatusNotFound,
		},
		{
			name:            "should get error if units is invalid",
			ingredientIDstr: "1",
			payload: `{
				"units": 0,
				"price": 12.5
			}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"units should be more than 0"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:            "should get error if price is invalid",
			ingredientIDstr: "1",
			payload: `{
				"units": 5
			}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"price is invalid"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/ingredients/"+tc.ingredientIDstr+"/stock", bytes.NewBufferString(tc.payload))
			require.NoError(t, err)
			rr := makeRequest(t, clock, func(useCases *usecases.UseCases) error {
				useCases.Ingredients.Create(context.Background(), ingredients.CreateIngredientOptions{
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
