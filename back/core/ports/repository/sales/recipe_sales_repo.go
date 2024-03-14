package salesrepo

import (
	"context"
	"costly/core/model"
	"costly/core/ports/database"
)

type RecipeSalesRepository interface {
	Add(ctx context.Context, recipeSales *model.RecipeSales) error
}

type repository struct {
	db database.TX
}

func New(db database.TX) RecipeSalesRepository {
	return &repository{db}
}

func (r *repository) Add(ctx context.Context, recipeSales *model.RecipeSales) error {
	result, err := r.db.ExecContext(ctx, "INSERT INTO sold_recipes_history (recipe_id, units, created_at) VALUES (?, ?, ?)", recipeSales.RecipeID, recipeSales.Units, recipeSales.CreatedAt)
	if err != nil {
		return err
	}
	recipeSalesID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	recipeSales.ID = recipeSalesID
	return nil
}
