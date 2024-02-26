package repository

import (
	"context"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"database/sql"
	"fmt"
	"time"
)

type RecipeRepository interface {
	GetRecipe(ctx context.Context, id int64) (domain.Recipe, error)
	GetRecipes(ctx context.Context) ([]domain.Recipe, error)
	CreateRecipe(ctx context.Context, name string, ingredients []RecipeIngredientInput) (domain.Recipe, error)
}

type recipeRepository struct {
	db     *database.Database
	clock  clock.Clock
	logger logger.Logger
}

type RecipeIngredientInput struct {
	IngredientID int64 `json:"id"`
	Units        int   `json:"units"`
}

type recipeDB struct {
	id           int64
	name         string
	createdAt    time.Time
	lastModified time.Time
}

func NewRecipeRepository(db *database.Database, clock clock.Clock, logger logger.Logger) RecipeRepository {
	return &recipeRepository{db, clock, logger}
}

func (r *recipeRepository) GetRecipe(ctx context.Context, id int64) (domain.Recipe, error) {

	row := r.db.QueryRowContext(ctx, "SELECT * FROM recipe WHERE id = ?", id)
	var recipe recipeDB
	err := row.Scan(&recipe.id, &recipe.name, &recipe.createdAt, &recipe.lastModified)
	if err == sql.ErrNoRows {
		return domain.Recipe{}, ErrNotFound
	} else if err != nil {
		return domain.Recipe{}, err
	}

	rows, err := r.db.QueryContext(ctx, "SELECT i.*, ri.units FROM ingredient i JOIN recipe_ingredient ri ON i.id = ri.ingredient_id AND ri.recipe_id = ?", id)
	if err == sql.ErrNoRows {
		return domain.Recipe{}, ErrNotFound
	} else if err != nil {
		return domain.Recipe{}, err
	}

	recipeIngredients := []domain.RecipeIngredient{}
	for rows.Next() {
		var ingredient domain.Ingredient
		var units int
		err := rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified, &units)
		if err != nil {
			return domain.Recipe{}, err
		}
		recipeIngredients = append(recipeIngredients, domain.RecipeIngredient{
			Ingredient: ingredient,
			Units:      units,
		})
	}

	return domain.Recipe{
		ID:           recipe.id,
		Name:         recipe.name,
		Ingredients:  recipeIngredients,
		CreatedAt:    recipe.createdAt,
		LastModified: recipe.lastModified,
	}, nil
}

func (r *recipeRepository) GetRecipes(ctx context.Context) ([]domain.Recipe, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM recipe")
	if err != nil {
		return nil, err
	}

	recipesDB := []recipeDB{}
	for rows.Next() {
		var recipe recipeDB
		if err := rows.Scan(&recipe.id, &recipe.name, &recipe.createdAt, &recipe.lastModified); err != nil {
			return nil, err
		}
		recipesDB = append(recipesDB, recipe)
	}

	recipes := []domain.Recipe{}
	for _, recipeDB := range recipesDB {
		rows, err := r.db.QueryContext(ctx, "SELECT i.*, ri.units FROM ingredient i JOIN recipe_ingredient ri ON i.id = ri.ingredient_id AND ri.recipe_id = ?", recipeDB.id)
		if err != nil {
			return nil, err
		}
		recipeIngredients := []domain.RecipeIngredient{}
		for rows.Next() {
			var ingredient domain.Ingredient
			var units int
			err := rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified, &units)
			if err != nil {
				return nil, err
			}
			recipeIngredients = append(recipeIngredients, domain.RecipeIngredient{
				Ingredient: ingredient,
				Units:      units,
			})
		}
		recipes = append(recipes, domain.Recipe{
			ID:           recipeDB.id,
			Name:         recipeDB.name,
			Ingredients:  recipeIngredients,
			CreatedAt:    recipeDB.createdAt,
			LastModified: recipeDB.lastModified,
		})
	}
	return recipes, nil
}

func (r *recipeRepository) CreateRecipe(ctx context.Context, name string, ingredients []RecipeIngredientInput) (domain.Recipe, error) {
	now := r.clock.Now()
	var recipeID int64 = -1
	recipeIngredients := []domain.RecipeIngredient{}
	if err := r.db.WithTx(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, "INSERT INTO recipe (name, created_at, last_modified) VALUES (?, ?, ?)", name, now, now)
		if err != nil {
			return err
		}

		recipeID, err = result.LastInsertId()
		if err != nil {
			return err
		}

		for _, recipeIngredient := range ingredients {
			// check ingredient exists
			// TODO: Make this more performant instead of querying one by one and before creating recipe
			row := tx.QueryRowContext(ctx, "SELECT * FROM ingredient WHERE id = ?", recipeIngredient.IngredientID)
			var ingredient domain.Ingredient
			err := row.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified)
			if err == sql.ErrNoRows {
				// return fmt.Errorf("error creating recipe with unexistent ingredient")
				return ErrBadOpts
			} else if err != nil {
				return err
			}

			recipeIngredients = append(recipeIngredients, domain.RecipeIngredient{
				Ingredient: ingredient,
				Units:      recipeIngredient.Units,
			})

			_, err = tx.ExecContext(ctx, "INSERT INTO recipe_ingredient (recipe_id, ingredient_id, units) VALUES (?, ?, ?)", recipeID, ingredient.ID, recipeIngredient.Units)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return domain.Recipe{}, fmt.Errorf("failed to create recipe: %w", err)
	}

	return domain.Recipe{
		ID:           recipeID,
		Name:         name,
		Ingredients:  recipeIngredients,
		CreatedAt:    now,
		LastModified: now,
	}, nil
}
