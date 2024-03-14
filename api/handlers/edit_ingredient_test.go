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

func TestHandleEditIngredient(t *testing.T) {
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
			name:            "should edit ingredient if existent",
			ingredientIDstr: "1",
			payload: `{
				"name": "green tea",
				"unit": "gr",
				"price": 10.0
			}`,
			expected:   "",
			statusCode: http.StatusNoContent,
		},
		{
			name:            "should get error if editing unexistent ingredient",
			ingredientIDstr: "123",
			payload: `{
				"name": "green tea",
				"unit": "gr",
				"price": 10.0
			}`,
			expected:   "",
			statusCode: http.StatusNotFound,
		},
		{
			name:            "should get error if name is valid and unit and price are not present",
			ingredientIDstr: "1",
			payload: `{
				"name": "aValidName"
			}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"unit is invalid"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:            "should get error if name is valid and unit not valid",
			ingredientIDstr: "1",
			payload: `{
				"name": "aValidName",
				"unit": "invalidUnit",
				"price": 12.32
			}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"unit is invalid"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:            "should get error if name is empty",
			ingredientIDstr: "1",
			payload: `{
				"name": "",
				"unit": "gr",
				"price": 123.0
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
			name:            "should get error if name and name are valid, but price is 0",
			ingredientIDstr: "1",
			payload: `{
				"name": "aValidNamE",
				"unit": "gr",
				"price": 0
			}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"price is invalid"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:            "should get error if bad request id",
			ingredientIDstr: "badID",
			payload: `{
				"name": "green tea",
				"unit": "gr",
				"price": 10.0
			}`,
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
			req, err := http.NewRequest("PUT", "/ingredients/"+tc.ingredientIDstr, bytes.NewBufferString(tc.payload))
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
