package handlers_test

import (
	"context"
	"costly/api/handlers"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func runGetIngredientHandler(t *testing.T, clock clock.Clock, ingredientIDstr string) *httptest.ResponseRecorder {
	logger, _ := logger.NewLogger("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	repo := repository.NewIngredientRepository(db, clock, logger)
	repo.CreateIngredient(context.Background(), repository.CreateIngredientOptions{
		Name:  "ingredientName",
		Price: 12.43,
		Unit:  domain.Gram,
	})
	handler := handlers.GetIngredientHandler(repo)

	req, err := http.NewRequest("GET", "/ingredients/"+ingredientIDstr, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/ingredients/{ingredientID}", handler)
	mux.ServeHTTP(rr, req)

	return rr
}

func TestHandleGetIngredient(t *testing.T) {
	clock := new(clockMock)
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
			rr := runGetIngredientHandler(t, clock, tc.ingredientIDstr)
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
