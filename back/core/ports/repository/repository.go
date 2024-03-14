package repo

import (
	"context"
	"costly/core/ports/database"
	ingredientrepo "costly/core/ports/repository/ingredient"
	reciperepo "costly/core/ports/repository/recipe"
	recipeviewrepo "costly/core/ports/repository/recipe_view"
	salesrepo "costly/core/ports/repository/sales"
	stockrepo "costly/core/ports/repository/stock"
)

type Repository interface {
	IngredientStocks() stockrepo.IngredientStockRepository
	Ingredients() ingredientrepo.IngredientRepository
	Recipes() reciperepo.RecipeRepository
	RecipeSales() salesrepo.RecipeSalesRepository
	RecipeViews() recipeviewrepo.RecipeViewRepository
	Atomic(ctx context.Context, fn func(repo Repository) error) error
}

type repository struct {
	db      database.Database
	session database.TX
}

func New(db database.Database) Repository {
	return &repository{db, db}
}

func (r *repository) IngredientStocks() stockrepo.IngredientStockRepository {
	return stockrepo.New(r.session)
}

func (r *repository) Ingredients() ingredientrepo.IngredientRepository {
	return ingredientrepo.New(r.session)
}

func (r *repository) Recipes() reciperepo.RecipeRepository {
	return reciperepo.New(r.session)
}

func (r *repository) RecipeSales() salesrepo.RecipeSalesRepository {
	return salesrepo.New(r.session)
}

func (r *repository) RecipeViews() recipeviewrepo.RecipeViewRepository {
	return recipeviewrepo.New(r.session)
}

func (r *repository) Atomic(ctx context.Context, fn func(repo Repository) error) (err error) {
	return r.db.WithTx(ctx, func(tx database.TX) error {
		newRepo := &repository{
			db: r.db,
			// injecting the new trx handle
			session: tx,
		}
		return fn(newRepo)
	})
}
