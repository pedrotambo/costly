package handlers_test

import (
	"bytes"
	"costly/api/handlers"
	"costly/core/components/clock"
	"costly/core/components/database"
	"costly/core/components/ingredients"
	"costly/core/components/logger"
	"costly/core/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func runCreateIngredientHandler(t *testing.T, clock clock.Clock, reqBody io.Reader) *httptest.ResponseRecorder {
	logger, _ := logger.New("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientComponent := ingredients.New(db, clock, logger)
	handler := handlers.CreateIngredientHandler(ingredientComponent)

	req, err := http.NewRequest("POST", "/ingredients", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}

func TestHandleCreateIngredient(t *testing.T) {
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
				"name": "recipeName",
				"price": 12.43,
				"unit": "gr"
			}`,
			expected: `{
				"id":1,
				"name":"recipeName",
				"unit":"gr",
				"price":12.43,
				"units_in_stock":0,
				"created_at":"1970-01-01T00:00:12.345Z",
				"last_modified":"1970-01-01T00:00:12.345Z"
			}`,
			statusCode: http.StatusCreated,
		},
		{
			name:    "should return error if unit is invalid",
			payload: `{"name": "validName", "price": 12.43, "unit": "notAtGr"}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"unit is invalid"
				}
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:    "should return error if name is invalid",
			payload: `{"name": "", "price": 12.43, "unit": "gr"}`,
			expected: `{
				"error": {
					"code":"INVALID_INPUT",
					"message":"name is invalid"
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
			rr := runCreateIngredientHandler(t, clock, bytes.NewBufferString(tc.payload))
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
