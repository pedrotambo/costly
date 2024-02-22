package costs_test

import (
	"context"
	"costly/core/domain"
	costs "costly/core/logic"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type recipeRepositoryMock struct {
	mock.Mock
}

func (m *recipeRepositoryMock) Get(ctx context.Context, recipeId string) (*domain.Recipe, error) {
	args := m.Called(recipeId)
	value := args.Get(0)
	recipe, ok := value.(*domain.Recipe)
	if !ok {
		return &domain.Recipe{}, fmt.Errorf("Error getting recipe")
	}
	return recipe, nil
}

func (m *recipeRepositoryMock) Save(ctx context.Context, recipe *domain.Recipe) error {
	args := m.Called(recipe)
	return args.Error(0)
}

func TestRecipeCost(t *testing.T) {

	egg := domain.Ingredient{
		Name:  "Egg",
		Price: 5.0,
		Unit:  domain.Units,
	}

	salt := domain.Ingredient{
		Name:  "Salt",
		Price: 1.0,
		Unit:  domain.Gram,
	}

	meat := domain.Ingredient{
		Name:  "Meat",
		Price: 10.0,
		Unit:  domain.Gram,
	}

	// t.Run("cost of a recipe is the sum of its ingredients and units of them", func(t *testing.T) {

	// 	recipeId := "someId"
	// 	recipeRepositoryMock := new(recipeRepositoryMock)
	// 	recipeRepositoryMock.On("Get", recipeId).Return(&domain.Recipe{
	// 		Name: "someRecipe",
	// 		Ingredients: []domain.RecipeIngredient{
	// 			{Ingredient: egg, Units: 5},
	// 			{Ingredient: salt, Units: 10},
	// 			{Ingredient: meat, Units: 500},
	// 		},
	// 	}, nil)

	// 	costService := costs.NewCostService(recipeRepositoryMock)

	// 	cost, _ := costService.GetRecipeCost(context.Background(), recipeId)

	// 	assert.Equal(t, 25.0+10.0+5000.0, cost)
	// })

	t.Run("cost of a recipe is the sum of its ingredients and units of them", func(t *testing.T) {

		costService := costs.NewCostService()
		cost := costService.GetRecipeCost(context.Background(), &domain.Recipe{
			Name: "someRecipe",
			Ingredients: []domain.RecipeIngredient{
				{Ingredient: egg, Units: 5},
				{Ingredient: salt, Units: 10},
				{Ingredient: meat, Units: 500},
			},
		})

		assert.Equal(t, 25.0+10.0+5000.0, cost)
	})
}
