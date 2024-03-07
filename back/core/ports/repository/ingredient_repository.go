package repository

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
}

type CreateIngredientOptions struct {
	Name  string
	Price float64
	Unit  domain.Unit
}

type ingredientRepository struct {
	db     *database.Database
	clock  clock.Clock
	logger logger.Logger
}

func NewIngredientRepository(db *database.Database, clock clock.Clock, logger logger.Logger) IngredientRepository {
	return &ingredientRepository{db, clock, logger}
}

func (r *ingredientRepository) GetIngredient(ctx context.Context, id int64) (domain.Ingredient, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM ingredient WHERE id = ?", id)

	var ingredient domain.Ingredient
	err := row.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified)
	if err == sql.ErrNoRows {
		return domain.Ingredient{}, ErrNotFound
	} else if err != nil {
		return domain.Ingredient{}, err
	}
	return ingredient, nil
}

func (r *ingredientRepository) GetIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM ingredient")
	if err != nil {
		return nil, err
	}

	ingredients := []domain.Ingredient{}
	for rows.Next() {
		var ingredient domain.Ingredient
		if err := rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return ingredients, nil
}

func (r *ingredientRepository) CreateIngredient(ctx context.Context, ingredientOpts CreateIngredientOptions) (domain.Ingredient, error) {
	now := r.clock.Now()
	var ingredientID int64 = -1
	result, err := r.db.ExecContext(ctx, "INSERT INTO ingredient (name, unit, price, created_at, last_modified) VALUES (?, ?, ?, ?, ?)", ingredientOpts.Name, ingredientOpts.Unit, ingredientOpts.Price, now, now)
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
		CreatedAt:    now,
		LastModified: now,
	}, nil
}

func (r *ingredientRepository) EditIngredient(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) (domain.Ingredient, error) {
	now := r.clock.Now()
	row := r.db.QueryRowContext(ctx, "UPDATE ingredient SET name = ?, unit = ?, price = ?, last_modified = ? WHERE id = ? RETURNING *",
		ingredientOpts.Name, ingredientOpts.Unit, ingredientOpts.Price, now, ingredientID)

	var ingredient domain.Ingredient
	err := row.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified)
	if err == sql.ErrNoRows {
		r.logger.Error(ErrNotFound, "error updating unexistent ingredient")
		return domain.Ingredient{}, ErrNotFound
	} else if err != nil {
		r.logger.Error(err, "error updating ingredient")
		return domain.Ingredient{}, err
	}
	return ingredient, nil
}
