package handlers_test

import (
	"bytes"
	"context"
	"costly/api/handlers"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/usecases"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func runUpdateIngredientStockHandler(t *testing.T, clock clock.Clock, ingredientIDstr string, reqBody io.Reader) *httptest.ResponseRecorder {
	logger, _ := logger.New("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	repo := rpst.NewIngredientRepository(db, clock, logger)
	ingredientUsecases := usecases.NewIngredientUseCases(repo, clock)
	ingredientUsecases.CreateIngredient(context.Background(), usecases.CreateIngredientOptions{
		Name:  "ingredientName",
		Price: 12.43,
		Unit:  model.Gram,
	})
	handler := handlers.UpdateIngredientStockHandler(repo)

	req, err := http.NewRequest("PUT", "/ingredients/stock/"+ingredientIDstr, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/ingredients/stock/{ingredientID}", handler)
	mux.ServeHTTP(rr, req)

	return rr
}

func TestHandleUpdateIngredientStock(t *testing.T) {
	clock := new(clockMock)
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
			name:            "should update ingredient stock if existent",
			ingredientIDstr: "1",
			payload: `{
				"new_units": 5,
				"price": 12.5
			}`,
			expected:   "",
			statusCode: http.StatusNoContent,
		},
		{
			name:            "should get error if updating stock of unexistent ingredient",
			ingredientIDstr: "123",
			payload: `{
				"new_units": 5,
				"price": 12.5
			}`,
			expected:   "",
			statusCode: http.StatusNotFound,
		},
		{
			name:            "should get error if new_units is invalid",
			ingredientIDstr: "1",
			payload: `{
				"new_units": 0,
				"price": 12.5
			}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"new_units should be more than 0"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:            "should get error if price is invalid",
			ingredientIDstr: "1",
			payload: `{
				"new_units": 5
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
			rr := runUpdateIngredientStockHandler(t, clock, tc.ingredientIDstr, bytes.NewBufferString(tc.payload))
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
