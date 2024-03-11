package handlers_test

import (
	"bytes"
	"context"
	"costly/api/handlers"
	"costly/core/components/clock"
	"costly/core/components/database"
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/mocks"
	"costly/core/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func runAddIngredientStockHandler(t *testing.T, clock clock.Clock, ingredientIDstr string, reqBody io.Reader) *httptest.ResponseRecorder {
	logger, _ := logger.New("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientComponent := ingredients.New(db, clock, logger)
	ingredientComponent.CreateIngredient(context.Background(), ingredients.CreateIngredientOptions{
		Name:  "ingredientName",
		Price: 12.43,
		Unit:  model.Gram,
	})
	handler := handlers.AddIngredientStockHandler(ingredientComponent)

	req, err := http.NewRequest("PUT", "/ingredients/"+ingredientIDstr+"/stock", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/ingredients/{ingredientID}/stock", handler)
	mux.ServeHTTP(rr, req)

	return rr
}

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
			rr := runAddIngredientStockHandler(t, clock, tc.ingredientIDstr, bytes.NewBufferString(tc.payload))
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
