package rpst

import (
	"context"
	"costly/core/components/database"
	"costly/core/model"
)

type rowMapper[T any] func(rowScanner database.RowScanner) (T, error)

func queryRowAndMap[T any](ctx context.Context, db database.RowQuerier, rowMapper rowMapper[T], query string, args ...any) (T, error) {
	row := db.QueryRowContext(ctx, query, args...)
	return rowMapper(row)
}

func queryAndMap[T any](ctx context.Context, db database.RowsQuerier, rowMapper rowMapper[T], query string, args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	ts := []T{}
	for rows.Next() {
		t, err := rowMapper(rows)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	return ts, nil
}

func mapToIngredient(rowScanner database.RowScanner) (model.Ingredient, error) {
	var ingredient model.Ingredient
	err := rowScanner.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified, &ingredient.UnitsInStock)
	return ingredient, err
}

func mapToRecipeIngredient(rowScanner database.RowScanner) (model.RecipeIngredient, error) {
	var ingredient model.Ingredient
	var recipeUnits int
	err := rowScanner.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified, &ingredient.UnitsInStock, &recipeUnits)
	return model.RecipeIngredient{
		Ingredient: ingredient,
		Units:      recipeUnits,
	}, err
}

func mapToRecipeDB(rowScanner database.RowScanner) (recipeDB, error) {
	var recipe recipeDB
	err := rowScanner.Scan(&recipe.id, &recipe.name, &recipe.createdAt, &recipe.lastModified)
	return recipe, err
}

func mapToIngredientStock(rowScanner database.RowScanner) (model.IngredientStock, error) {
	var stock model.IngredientStock
	err := rowScanner.Scan(&stock.ID, &stock.IngredientID, &stock.Units, &stock.Price, &stock.CreatedAt)
	return stock, err
}
