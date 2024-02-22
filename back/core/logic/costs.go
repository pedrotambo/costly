package costs

import (
	"context"
	"costly/core/domain"
)

type CostService interface {
	// GetRecipeCost(ctx context.Context, recipeId int64) (float64, error)
	GetRecipeCost(ctx context.Context, recipe *domain.Recipe) float64
	// GetRecipeCost(ctx context.Context, recipe *domain.Recipe) float64
}

type costService struct {
	// recipeRepository repository.RecipeRepository
}

// func NewCostService(RecipeRepository repository.RecipeRepository) *costService {
func NewCostService() *costService {
	return &costService{
		// RecipeRepository,
	}
}

// func (cs *costService) GetRecipeCost(ctx context.Context, recipeId int64) (float64, error) {
// 	recipe, err := cs.recipeRepository.GetRecipe(ctx, recipeId)

// 	if err != nil {
// 		return 0.0, err
// 	}

// 	cost := 0.0

// 	for _, ingredient := range recipe.Ingredients {
// 		cost += ingredient.Ingredient.Price * float64(ingredient.Units)
// 	}

// 	return cost, nil
// }

func (cs *costService) GetRecipeCost(ctx context.Context, recipe *domain.Recipe) float64 {
	cost := 0.0

	for _, ingredient := range recipe.Ingredients {
		cost += ingredient.Ingredient.Price * float64(ingredient.Units)
	}

	return cost
}
