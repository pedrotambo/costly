package rpst

import (
	"context"
	"costly/core/components/clock"
	"costly/core/components/database"
	"costly/core/components/logger"
	"costly/core/errs"
	"costly/core/model"
	"database/sql"
)

type IngredientGetter interface {
	GetIngredient(ctx context.Context, id int64) (model.Ingredient, error)
}

type IngredientsGetter interface {
	GetIngredients(ctx context.Context) ([]model.Ingredient, error)
}

type IngredientRepository interface {
	SaveIngredient(ctx context.Context, ingredient *model.Ingredient) error
	UpdateIngredient(ctx context.Context, ingredientID int64, updateFunc func(ingredient *model.Ingredient) error) error
	IngredientGetter
	IngredientsGetter
	SaveIngredientStock(ctx context.Context, ingredientStock *model.IngredientStock) error
	GetIngredientStockHistory(ctx context.Context, ingredientID int64) ([]model.IngredientStock, error)
}

type NewStockOptions struct {
	NewUnits int `json:"new_units"`
	Price    float64
}

type ingredientRepository struct {
	db     database.Database
	clock  clock.Clock
	logger logger.Logger
}

func New(db database.Database, clock clock.Clock, logger logger.Logger) IngredientRepository {
	return &ingredientRepository{db, clock, logger}
}

func (r *ingredientRepository) GetIngredient(ctx context.Context, id int64) (model.Ingredient, error) {
	ingredient, err := queryRowAndMap(ctx, r.db, mapToIngredient, "SELECT * FROM ingredient WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return model.Ingredient{}, errs.ErrNotFound
	} else if err != nil {
		return model.Ingredient{}, err
	}
	return ingredient, nil
}

func (r *ingredientRepository) GetIngredients(ctx context.Context) ([]model.Ingredient, error) {
	ingredients, err := queryAndMap(ctx, r.db, mapToIngredient, "SELECT * FROM ingredient")
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (r *ingredientRepository) SaveIngredient(ctx context.Context, ingredient *model.Ingredient) error {
	result, err := r.db.ExecContext(ctx, "INSERT INTO ingredient (name, unit, price, units_in_stock, created_at, last_modified) VALUES (?, ?, ?, ?, ?, ?)",
		ingredient.Name, ingredient.Unit, ingredient.Price, ingredient.UnitsInStock, ingredient.CreatedAt, ingredient.LastModified)

	if err != nil {
		return err
	}

	ingredientID, err := result.LastInsertId()

	if err != nil {
		return err
	}
	ingredient.ID = ingredientID
	return nil
}

func (r *ingredientRepository) UpdateStock(ctx context.Context, ingredientID int64, newStockOptions NewStockOptions) (model.Ingredient, error) {
	now := r.clock.Now()
	var ingredient model.Ingredient
	if err := r.db.WithTx(ctx, func(tx database.TX) error {
		updatedIngredient, err := queryRowAndMap(ctx, tx, mapToIngredient, "UPDATE ingredient SET units_in_stock = units_in_stock + ?, price = ?, last_modified = ? WHERE id = ? RETURNING *", newStockOptions.NewUnits, newStockOptions.Price, now, ingredientID)
		if err == sql.ErrNoRows {
			r.logger.Error(errs.ErrNotFound, "error updating unexistent ingredient")
			return errs.ErrNotFound
		} else if err != nil {
			r.logger.Error(err, "error updating ingredient")
			return err
		}
		ingredient = updatedIngredient
		_, err = tx.ExecContext(ctx, "INSERT INTO stock_history (ingredient_id, units, price, created_at) VALUES (?, ?, ?, ?)", ingredientID, newStockOptions.NewUnits, newStockOptions.Price, now)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return model.Ingredient{}, err
	}
	return ingredient, nil
}

func (r *ingredientRepository) SaveIngredientStock(ctx context.Context, ingredientStock *model.IngredientStock) error {
	return r.db.WithTx(ctx, func(tx database.TX) error {
		res, err := tx.ExecContext(ctx, "UPDATE ingredient SET units_in_stock = units_in_stock + ?, price = ?, last_modified = ? WHERE id = ?", ingredientStock.Units, ingredientStock.Price, ingredientStock.CreatedAt, ingredientStock.IngredientID)

		if err != nil {
			r.logger.Error(err, "error updating ingredient stock")
			return err
		}

		if rows, err := res.RowsAffected(); err != nil {
			r.logger.Error(err, "error updating ingredient stock: "+err.Error())
			return err
		} else if rows == 0 {
			r.logger.Error(err, "error updating ingredient stock: inexisten ingredient")
			return errs.ErrNotFound
		}

		result, err := tx.ExecContext(ctx, "INSERT INTO stock_history (ingredient_id, units, price, created_at) VALUES (?, ?, ?, ?)", ingredientStock.IngredientID, ingredientStock.Units, ingredientStock.Price, ingredientStock.CreatedAt)

		if err != nil {
			r.logger.Error(err, "error saving ingredient stock")
			return err
		}

		ingredientStockID, err := result.LastInsertId()

		if err != nil {
			return err
		}

		ingredientStock.ID = ingredientStockID
		return nil
	})
}

func (r *ingredientRepository) GetIngredientStockHistory(ctx context.Context, ingredientID int64) ([]model.IngredientStock, error) {
	ingredients, err := queryAndMap(ctx, r.db, mapToIngredientStock, "SELECT * FROM stock_history")
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (r *ingredientRepository) UpdateIngredient(ctx context.Context, ingredientID int64, updateFunc func(ingredient *model.Ingredient) error) error {
	ingredient, err := r.GetIngredient(ctx, ingredientID)
	if err != nil {
		return err
	}
	updateFunc(&ingredient)
	_, err = queryRowAndMap(ctx, r.db, mapToIngredient, "UPDATE ingredient SET name = ?, unit = ?, price = ?, last_modified = ? WHERE id = ? RETURNING *", ingredient.Name, ingredient.Unit, ingredient.Price, ingredient.LastModified, ingredient.ID)
	if err == sql.ErrNoRows {
		return errs.ErrNotFound
	} else if err != nil {
		r.logger.Error(err, "error updating ingredient")
		return err
	}
	return nil
}
