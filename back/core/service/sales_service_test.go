package service_test

import (
	"context"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"costly/core/ports/rpst"
	"costly/core/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterSoldRecipes(t *testing.T) {
	clock := clock.New()
	logger, _ := logger.New("debug")
	db, _ := database.NewFromDatasource(":memory:", logger)
	repo := rpst.New(db, clock, logger)

	salesService := service.NewSalesService(repo)

	// err := salesService.RegisterSoldRecipes(context.Background(), service.SoldRecipes{RecipeID: 1, SoldUnits: 5})

	// require.Error(t, err)

	t.Run("should return error if sold recipes correspond to unexisten recipe", func(t *testing.T) {
		err := salesService.RegisterSoldRecipes(context.Background(), service.SoldRecipes{RecipeID: 1, SoldUnits: 5})
		require.Error(t, err)
		assert.Equal(t, rpst.ErrNotFound, err)
	})

	t.Run("should ASDF", func(t *testing.T) {
		ctx := context.Background()
		salt, _ := repo.CreateIngredient(ctx, rpst.CreateIngredientOptions{
			Name:  "salt",
			Unit:  domain.Gram,
			Price: 1.0,
		})
		meat, _ := repo.CreateIngredient(ctx, rpst.CreateIngredientOptions{
			Name:  "meat",
			Unit:  domain.Gram,
			Price: 10.0,
		})

		recipe, _ := repo.CreateRecipe(ctx, rpst.CreateRecipeOptions{
			Name: "meat with salt",
			Ingredients: []rpst.RecipeIngredientOptions{
				{
					ID:    meat.ID,
					Units: 500,
				},
				{
					ID:    salt.ID,
					Units: 5,
				},
			},
		})

		repo.UpdateStock(ctx, salt.ID, rpst.NewStockOptions{
			NewUnits: 100,
			Price:    1.0,
		})

		repo.UpdateStock(ctx, meat.ID, rpst.NewStockOptions{
			NewUnits: 1000,
			Price:    10.0,
		})

		numberOfSoldRecipes := 5
		err := salesService.RegisterSoldRecipes(context.Background(), service.SoldRecipes{RecipeID: recipe.ID, SoldUnits: numberOfSoldRecipes})
		require.NoError(t, err)

		// modifiedSalt, _ := repo.GetIngredient(ctx, salt.ID)
		// expectedSaltStock := 100.0 - numberOfSoldRecipes*5
		// assert.Equal(t, expectedSaltStock, modifiedSalt.UnitsInStock)

	})
}
