package stockrepo

import (
	"context"
	"costly/core/errs"
	"costly/core/model"
	"costly/core/ports/database"
	"database/sql"

	"github.com/mattn/go-sqlite3"
)

type IngredientStockRepository interface {
	Add(ctx context.Context, ingredientStock *model.IngredientStock) error
	Find(ctx context.Context, ingredientStockID int64) (model.IngredientStock, error)
}

type ingredientStockRepository struct {
	db database.TX
}

func New(db database.TX) IngredientStockRepository {
	return &ingredientStockRepository{db}
}

func (r *ingredientStockRepository) Add(ctx context.Context, ingredientStock *model.IngredientStock) error {
	result, err := r.db.ExecContext(ctx, "INSERT INTO stock_history (ingredient_id, units, price, created_at) VALUES (?, ?, ?, ?)", ingredientStock.IngredientID, ingredientStock.Units, ingredientStock.Price, ingredientStock.CreatedAt)
	if err != nil {
		if sqlError, ok := err.(sqlite3.Error); ok {
			if sqlError.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return errs.ErrNotFound
			}
		}
		return err
	}
	ingredientStockID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ingredientStock.ID = ingredientStockID
	return nil
}

func (r *ingredientStockRepository) Find(ctx context.Context, ingredientStockID int64) (model.IngredientStock, error) {
	stock, err := database.QueryRowAndMap(ctx, r.db, mapToIngredientStock, "SELECT * FROM stock_history WHERE id = ?", ingredientStockID)
	if err == sql.ErrNoRows {
		return model.IngredientStock{}, errs.ErrNotFound
	} else if err != nil {
		return model.IngredientStock{}, err
	}
	return stock, nil
}

func mapToIngredientStock(rowScanner database.RowScanner) (model.IngredientStock, error) {
	var ingredientStock model.IngredientStock
	err := rowScanner.Scan(&ingredientStock.ID, &ingredientStock.IngredientID, &ingredientStock.Units, &ingredientStock.Price, &ingredientStock.CreatedAt)
	return ingredientStock, err
}
