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

func runGetIngredientsHandler(t *testing.T, clock clock.Clock, ingrOpts []repository.CreateIngredientOptions) *httptest.ResponseRecorder {
	logger, _ := logger.NewLogger("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	repo := repository.NewIngredientRepository(db, clock, logger)
	for _, opts := range ingrOpts {
		repo.CreateIngredient(context.Background(), opts)
	}

	handler := handlers.GetIngredientsHandler(repo)

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
	clock := new(clockMock)
	now := time.UnixMilli(12345).UTC()
	clock.On("Now").Return(now)

	testCases := []struct {
		name        string
		ingredients []repository.CreateIngredientOptions
		expected    string
		statusCode  int
	}{
		{
			name: "should get ingredients",
			ingredients: []repository.CreateIngredientOptions{
				{
					Name:  "ingr1",
					Price: 1.5,
					Unit:  domain.Gram,
				},
				{
					Name:  "ingr2",
					Price: 2.5,
					Unit:  domain.Gram,
				},
			},
			expected: `[
				{
					"id": 1,
					"name": "ingr1",
					"unit": "gr",
					"price": 1.5,
					"created_at": "1970-01-01T00:00:12.345Z",
					"last_modified": "1970-01-01T00:00:12.345Z"
				},
				{
					"id": 2,
					"name": "ingr2",
					"unit": "gr",
					"price": 2.5,
					"created_at": "1970-01-01T00:00:12.345Z",
					"last_modified": "1970-01-01T00:00:12.345Z"
				}
			]`,
			statusCode: http.StatusOK,
		},
		{
			name:        "should get empty ingredients",
			ingredients: []repository.CreateIngredientOptions{},
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
