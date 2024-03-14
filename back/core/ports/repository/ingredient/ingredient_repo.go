package ingredientrepo

import (
	"context"
	"costly/core/errs"
	"costly/core/model"
	"costly/core/ports/database"
	"database/sql"
	"time"
)

type IngredientRepository interface {
	Add(ctx context.Context, ingredient *model.Ingredient) error
	Update(ctx context.Context, ingredientID int64, updateFunc func(ingredient *model.Ingredient) error) error
	Find(ctx context.Context, id int64) (model.Ingredient, error)
	FindAll(ctx context.Context) ([]model.Ingredient, error)
	IncreaseStockAndUpdatePrice(ctx context.Context, ingredientID int64, units int, price float64, now time.Time) error
	DecreaseStock(ctx context.Context, ingredientID int64, unitsToDecrease int, now time.Time) error
}

type ingredientRepository struct {
	db database.TX
}

func New(db database.TX) IngredientRepository {
	return &ingredientRepository{db}
}

func (r *ingredientRepository) Find(ctx context.Context, id int64) (model.Ingredient, error) {
	ingredient, err := database.QueryRowAndMap(ctx, r.db, mapToIngredient, "SELECT * FROM ingredient WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return model.Ingredient{}, errs.ErrNotFound
	} else if err != nil {
		return model.Ingredient{}, err
	}
	return ingredient, nil
}

func (r *ingredientRepository) FindAll(ctx context.Context) ([]model.Ingredient, error) {
	ingredients, err := database.QueryAndMap(ctx, r.db, mapToIngredient, "SELECT * FROM ingredient")
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (r *ingredientRepository) Add(ctx context.Context, ingredient *model.Ingredient) error {
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

func (r *ingredientRepository) Update(ctx context.Context, ingredientID int64, updateFunc func(ingredient *model.Ingredient) error) error {
	ingredient, err := r.Find(ctx, ingredientID)
	if err != nil {
		return err
	}
	updateFunc(&ingredient)
	_, err = database.QueryRowAndMap(ctx, r.db, mapToIngredient, "UPDATE ingredient SET name = ?, unit = ?, price = ?, units_in_stock = ?, last_modified = ? WHERE id = ? RETURNING *",
		ingredient.Name, ingredient.Unit, ingredient.Price, ingredient.UnitsInStock, ingredient.LastModified, ingredient.ID)
	if err == sql.ErrNoRows {
		return errs.ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}

func (r *ingredientRepository) IncreaseStockAndUpdatePrice(ctx context.Context, ingredientID int64, units int, price float64, now time.Time) error {
	result, err := r.db.ExecContext(ctx, "UPDATE ingredient SET units_in_stock = units_in_stock + ?, price = ?, last_modified = ? WHERE id = ?", units, price, now, ingredientID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func (r *ingredientRepository) DecreaseStock(ctx context.Context, ingredientID int64, unitsToDecrease int, timeOfDecrease time.Time) error {
	result, err := r.db.ExecContext(ctx, "UPDATE ingredient SET units_in_stock = units_in_stock - ?, last_modified = ? WHERE id = ?", unitsToDecrease, timeOfDecrease, ingredientID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func mapToIngredient(rowScanner database.RowScanner) (model.Ingredient, error) {
	var ingredient model.Ingredient
	err := rowScanner.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified, &ingredient.UnitsInStock)
	return ingredient, err
}
