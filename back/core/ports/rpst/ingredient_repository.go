package rpst

import (
	"context"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"database/sql"
)

type IngredientRepository interface {
	GetIngredient(ctx context.Context, id int64) (domain.Ingredient, error)
	GetIngredients(ctx context.Context) ([]domain.Ingredient, error)
	CreateIngredient(ctx context.Context, ingredientOpts CreateIngredientOptions) (domain.Ingredient, error)
	EditIngredient(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) (domain.Ingredient, error)
	UpdateStock(ctx context.Context, ingredientID int64, newStockOptions NewStockOptions) (domain.Ingredient, error)
	// ReduceStockUnits(ctx context.Context, unitsUsedByID map[int64]int) error
}

type CreateIngredientOptions struct {
	Name  string
	Price float64
	Unit  domain.Unit
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

func NewIngredientRepository(db database.Database, clock clock.Clock, logger logger.Logger) IngredientRepository {
	return &ingredientRepository{db, clock, logger}
}

func (r *ingredientRepository) GetIngredient(ctx context.Context, id int64) (domain.Ingredient, error) {
	ingredient, err := queryRowAndMap(ctx, r.db, mapToIngredient, "SELECT * FROM ingredient WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return domain.Ingredient{}, ErrNotFound
	} else if err != nil {
		return domain.Ingredient{}, err
	}
	return ingredient, nil
}

func (r *ingredientRepository) GetIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	ingredients, err := queryAndMap(ctx, r.db, mapToIngredient, "SELECT * FROM ingredient")
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (r *ingredientRepository) CreateIngredient(ctx context.Context, ingredientOpts CreateIngredientOptions) (domain.Ingredient, error) {
	now := r.clock.Now()
	var ingredientID int64 = -1
	result, err := r.db.ExecContext(ctx, "INSERT INTO ingredient (name, unit, price, units_in_stock, created_at, last_modified) VALUES (?, ?, ?, ?, ?, ?)", ingredientOpts.Name, ingredientOpts.Unit, ingredientOpts.Price, 0, now, now)
	if err != nil {
		return domain.Ingredient{}, err
	}
	ingredientID, err = result.LastInsertId()
	if err != nil {
		return domain.Ingredient{}, err
	}
	return domain.Ingredient{
		ID:           ingredientID,
		Name:         ingredientOpts.Name,
		Price:        ingredientOpts.Price,
		Unit:         ingredientOpts.Unit,
		UnitsInStock: 0,
		CreatedAt:    now,
		LastModified: now,
	}, nil
}

func (r *ingredientRepository) EditIngredient(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) (domain.Ingredient, error) {
	now := r.clock.Now()
	ingredient, err := queryRowAndMap(ctx, r.db, mapToIngredient, "UPDATE ingredient SET name = ?, unit = ?, price = ?, last_modified = ? WHERE id = ? RETURNING *", ingredientOpts.Name, ingredientOpts.Unit, ingredientOpts.Price, now, ingredientID)
	if err == sql.ErrNoRows {
		r.logger.Error(ErrNotFound, "error updating unexistent ingredient")
		return domain.Ingredient{}, ErrNotFound
	} else if err != nil {
		r.logger.Error(err, "error updating ingredient")
		return domain.Ingredient{}, err
	}
	return ingredient, nil
}

func (r *ingredientRepository) UpdateStock(ctx context.Context, ingredientID int64, newStockOptions NewStockOptions) (domain.Ingredient, error) {
	now := r.clock.Now()
	var ingredient domain.Ingredient
	if err := r.db.WithTx(ctx, func(tx database.TX) error {
		updatedIngredient, err := queryRowAndMap(ctx, tx, mapToIngredient, "UPDATE ingredient SET units_in_stock = units_in_stock + ?, price = ?, last_modified = ? WHERE id = ? RETURNING *", newStockOptions.NewUnits, newStockOptions.Price, now, ingredientID)
		if err == sql.ErrNoRows {
			r.logger.Error(ErrNotFound, "error updating unexistent ingredient")
			return ErrNotFound
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
		return domain.Ingredient{}, err
	}
	return ingredient, nil
}

// func (r *ingredientRepository) ReduceStockUnits(ctx context.Context, unitsUsedByID map[int64]int) error {
// 	now := r.clock.Now()
// 	var ingredient domain.Ingredient
// 	if err := r.db.WithTx(ctx, func(tx database.TX) error {
// 		for id, usedUnits := range unitsUsedByID {

// 		}
// 		updatedIngredient, err := queryRowAndMap(ctx, tx, mapToIngredient, "UPDATE ingredient SET units_in_stock = units_in_stock + ?, price = ?, last_modified = ? WHERE id = ? RETURNING *", newStockOptions.NewUnits, newStockOptions.Price, now, ingredientID)
// 		if err == sql.ErrNoRows {
// 			r.logger.Error(ErrNotFound, "error updating unexistent ingredient")
// 			return ErrNotFound
// 		} else if err != nil {
// 			r.logger.Error(err, "error updating ingredient")
// 			return err
// 		}
// 		ingredient = updatedIngredient
// 		_, err = tx.ExecContext(ctx, "INSERT INTO stock_history (ingredient_id, units, price, created_at) VALUES (?, ?, ?, ?)", ingredientID, newStockOptions.NewUnits, newStockOptions.Price, now)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		return domain.Ingredient{}, err
// 	}
// 	return ingredient, nil
// }
