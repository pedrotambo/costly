package rpst

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
	CreateRecipe(ctx context.Context, recipeOpts CreateRecipeOptions) (domain.Recipe, error)
}

type recipeRepository struct {
	db     database.Database
	clock  clock.Clock
	logger logger.Logger
}

type RecipeIngredientOptions struct {
	ID    int64
	Units int
}

type CreateRecipeOptions struct {
	Name        string
	Ingredients []RecipeIngredientOptions
}

type recipeDB struct {
	id           int64
	name         string
	createdAt    time.Time
	lastModified time.Time
}

func NewRecipeRepository(db database.Database, clock clock.Clock, logger logger.Logger) RecipeRepository {
	return &recipeRepository{db, clock, logger}
}

func (r *recipeRepository) GetRecipe(ctx context.Context, id int64) (domain.Recipe, error) {
	recipeDB, err := queryRowAndMap(ctx, r.db, mapToRecipeDB, "SELECT * FROM recipe WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return domain.Recipe{}, ErrNotFound
	} else if err != nil {
		return domain.Recipe{}, err
	}
	recipeIngredients, err := queryAndMap(ctx, r.db, mapToRecipeIngredient, "SELECT i.*, ri.units FROM ingredient i JOIN recipe_ingredient ri ON i.id = ri.ingredient_id AND ri.recipe_id = ?", id)
	if err != nil {
		return domain.Recipe{}, err
	}
	return domain.Recipe{
		ID:           recipeDB.id,
		Name:         recipeDB.name,
		Ingredients:  recipeIngredients,
		CreatedAt:    recipeDB.createdAt,
		LastModified: recipeDB.lastModified,
	}, nil
}

func (r *recipeRepository) GetRecipes(ctx context.Context) ([]domain.Recipe, error) {
	recipesDB, err := queryAndMap(ctx, r.db, mapToRecipeDB, "SELECT * FROM recipe")
	if err != nil {
		return nil, err
	}
	recipes := []domain.Recipe{}
	for _, recipeDB := range recipesDB {
		recipeIngredients, err := queryAndMap(ctx, r.db, mapToRecipeIngredient, "SELECT i.*, ri.units FROM ingredient i JOIN recipe_ingredient ri ON i.id = ri.ingredient_id AND ri.recipe_id = ?", recipeDB.id)
		if err != nil {
			return nil, err
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

func (r *recipeRepository) CreateRecipe(ctx context.Context, recipeOpts CreateRecipeOptions) (domain.Recipe, error) {
	now := r.clock.Now()
	var recipeID int64 = -1
	recipeIngredients := []domain.RecipeIngredient{}
	if err := r.db.WithTx(ctx, func(tx database.TX) error {
		result, err := tx.ExecContext(ctx, "INSERT INTO recipe (name, created_at, last_modified) VALUES (?, ?, ?)", recipeOpts.Name, now, now)
		if err != nil {
			return err
		}

		recipeID, err = result.LastInsertId()
		if err != nil {
			return err
		}

		for _, recipeIngredient := range recipeOpts.Ingredients {
			// check ingredient exists
			// TODO: Make this more performant instead of querying one by one and before creating recipe
			ingredient, err := queryRowAndMap(ctx, tx, mapToIngredient, "SELECT * FROM ingredient WHERE id = ?", recipeIngredient.ID)
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
		Name:         recipeOpts.Name,
		Ingredients:  recipeIngredients,
		CreatedAt:    now,
		LastModified: now,
	}, nil
}
