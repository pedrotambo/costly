package handlers_test

import (
	"bytes"
	"costly/api/handlers"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/repository"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type clockMock struct {
	mock.Mock
}

func (m *clockMock) Now() time.Time {
	args := m.Called()
	value := args.Get(0)
	now, ok := value.(time.Time)
	if !ok {
		panic(fmt.Errorf("Error getting now"))
	}
	return now
}

func runCreatedIngredientHandler(t *testing.T, clock clock.Clock, reqBody io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest("POST", "/ingredients", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	logger, _ := logger.NewLogger("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	ingredientRepository := repository.NewIngredientRepository(db, clock, logger)

	createIngredientHandler := handlers.CreateIngredientHandler(ingredientRepository)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createIngredientHandler)
	handler.ServeHTTP(rr, req)

	return rr
}

func TestHandleCreateIngredient(t *testing.T) {
	clock := new(clockMock)
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
				"created_at":"1970-01-01T00:00:12.345Z",
				"last_modified":"1970-01-01T00:00:12.345Z"
			}`,
			statusCode: http.StatusCreated,
		},
		{
			name:    "should return error if unit is invalid",
			payload: `{"name": "validName", "price": 12.43, "unit": "notAtGr"}`,
			expected: `{
				"errors": [
					{
						"offending_field": "unit",
						"suggestion": "la unidad es inválida"
					}
				]
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:    "should return error if name is invalid",
			payload: `{"name": "", "price": 12.43, "unit": "gr"}`,
			expected: `{
				"errors": [
					{
						"offending_field": "name",
						"suggestion": "el name debe ser valido"
					}
				]
			}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "should return error if payload is invalid json",
			payload:    "invalid payload",
			expected:   "",
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := runCreatedIngredientHandler(t, clock, bytes.NewBufferString(tc.payload))
			assert.Equal(t, tc.statusCode, rr.Code)
			if tc.expected != rr.Body.String() {
				assert.JSONEq(t, tc.expected, rr.Body.String(), "Response body differs")
			}
		})
	}
}
