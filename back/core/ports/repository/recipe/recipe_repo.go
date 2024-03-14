package reciperepo

import (
	"context"
	"costly/core/errs"
	"costly/core/model"
	"costly/core/ports/database"
	"time"
)

type RecipeRepository interface {
	Add(ctx context.Context, recipe *model.Recipe) error
	Find(ctx context.Context, id int64) (model.Recipe, error)
	FindIngredients(ctx context.Context, recipeID int64) ([]model.RecipeIngredient, error)
}

type repository struct {
	db database.Database
}

func New(db database.Database) RecipeRepository {
	return &repository{db}
}

func (r *repository) Add(ctx context.Context, recipe *model.Recipe) error {
	return r.db.WithTx(ctx, func(tx database.Database) error {
		result, err := tx.ExecContext(ctx, "INSERT INTO recipe (name, created_at, last_modified) VALUES (?, ?, ?)", recipe.Name, recipe.CreatedAt, recipe.LastModified)
		if err != nil {
			return err
		}
		recipeID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for _, recipeIngredient := range recipe.Ingredients {
			_, err := tx.ExecContext(ctx, "INSERT INTO recipe_ingredient (recipe_id, ingredient_id, units) VALUES (?, ?, ?)", recipeID, recipeIngredient.ID, recipeIngredient.Units)
			if err != nil {
				return err
			}
		}

		recipe.ID = recipeID
		return nil
	})
}

func (r *repository) Find(ctx context.Context, id int64) (model.Recipe, error) {
	// This was carefully made to make only one query when selecting only one recipe.
	recipeWithIngredients, err := database.QueryAndMap(ctx, r.db, mapToRecipeWithIngredientsDB, "SELECT r.*, ri.ingredient_id, ri.units FROM recipe r JOIN recipe_ingredient ri ON r.id = ri.recipe_id WHERE r.id = ?", id)
	if err != nil {
		return model.Recipe{}, err
	}
	if len(recipeWithIngredients) == 0 {
		return model.Recipe{}, errs.ErrNotFound
	}
	recipeIngredients := []model.RecipeIngredient{}
	for _, ri := range recipeWithIngredients {
		recipeIngredients = append(recipeIngredients, model.RecipeIngredient{ID: ri.ingredientId, Units: ri.units})
	}
	return model.Recipe{
		ID:           recipeWithIngredients[0].id,
		Name:         recipeWithIngredients[0].name,
		Ingredients:  recipeIngredients,
		CreatedAt:    recipeWithIngredients[0].createdAt,
		LastModified: recipeWithIngredients[0].lastModified,
	}, nil
}

func (r *repository) FindIngredients(ctx context.Context, recipeID int64) ([]model.RecipeIngredient, error) {
	recipeIngredients, err := database.QueryAndMap(ctx, r.db, mapToRecipeIngredient, "SELECT ingredient_id, units FROM recipe_ingredient WHERE recipe_id = ?", recipeID)
	if err != nil {
		return nil, err
	}
	return recipeIngredients, nil
}

type recipeWithIngredient struct {
	id           int64
	name         string
	createdAt    time.Time
	lastModified time.Time
	ingredientId int64
	units        int
}

func mapToRecipeWithIngredientsDB(rowScanner database.RowScanner) (recipeWithIngredient, error) {
	var recipeWithIngredient recipeWithIngredient
	return recipeWithIngredient, rowScanner.Scan(&recipeWithIngredient.id, &recipeWithIngredient.name, &recipeWithIngredient.createdAt, &recipeWithIngredient.lastModified, &recipeWithIngredient.ingredientId, &recipeWithIngredient.units)
}

func mapToRecipeIngredient(rowScanner database.RowScanner) (model.RecipeIngredient, error) {
	var recipeIngredient model.RecipeIngredient
	return recipeIngredient, rowScanner.Scan(&recipeIngredient.ID, &recipeIngredient.Units)
}
