package rpst

import (
	"context"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
	"database/sql"
	"time"
)

type RecipeGetter interface {
	GetRecipe(ctx context.Context, id int64) (model.Recipe, error)
}

type RecipesGetter interface {
	GetRecipes(ctx context.Context) ([]model.Recipe, error)
}

type RecipeRepository interface {
	SaveRecipe(ctx context.Context, recipe *model.Recipe) error
	RecipeGetter
	RecipesGetter
	SaveRecipeSales(ctx context.Context, recipeSales *model.RecipeSales) error
}

type recipeRepository struct {
	db     database.Database
	clock  clock.Clock
	logger logger.Logger
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

func (r *recipeRepository) GetRecipe(ctx context.Context, id int64) (model.Recipe, error) {
	recipeDB, err := queryRowAndMap(ctx, r.db, mapToRecipeDB, "SELECT * FROM recipe WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return model.Recipe{}, ErrNotFound
	} else if err != nil {
		return model.Recipe{}, err
	}
	recipeIngredients, err := queryAndMap(ctx, r.db, mapToRecipeIngredient, "SELECT i.*, ri.units FROM ingredient i JOIN recipe_ingredient ri ON i.id = ri.ingredient_id AND ri.recipe_id = ?", id)
	if err != nil {
		return model.Recipe{}, err
	}
	return model.Recipe{
		ID:           recipeDB.id,
		Name:         recipeDB.name,
		Ingredients:  recipeIngredients,
		CreatedAt:    recipeDB.createdAt,
		LastModified: recipeDB.lastModified,
	}, nil
}

func getRecipe(ctx context.Context, tx database.TX, id int64) (model.Recipe, error) {
	recipeDB, err := queryRowAndMap(ctx, tx, mapToRecipeDB, "SELECT * FROM recipe WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return model.Recipe{}, ErrNotFound
	} else if err != nil {
		return model.Recipe{}, err
	}
	recipeIngredients, err := queryAndMap(ctx, tx, mapToRecipeIngredient, "SELECT i.*, ri.units FROM ingredient i JOIN recipe_ingredient ri ON i.id = ri.ingredient_id AND ri.recipe_id = ?", id)
	if err != nil {
		return model.Recipe{}, err
	}
	return model.Recipe{
		ID:           recipeDB.id,
		Name:         recipeDB.name,
		Ingredients:  recipeIngredients,
		CreatedAt:    recipeDB.createdAt,
		LastModified: recipeDB.lastModified,
	}, nil
}

func (r *recipeRepository) GetRecipes(ctx context.Context) ([]model.Recipe, error) {
	recipesDB, err := queryAndMap(ctx, r.db, mapToRecipeDB, "SELECT * FROM recipe")
	if err != nil {
		return nil, err
	}
	recipes := []model.Recipe{}
	for _, recipeDB := range recipesDB {
		recipeIngredients, err := queryAndMap(ctx, r.db, mapToRecipeIngredient, "SELECT i.*, ri.units FROM ingredient i JOIN recipe_ingredient ri ON i.id = ri.ingredient_id AND ri.recipe_id = ?", recipeDB.id)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, model.Recipe{
			ID:           recipeDB.id,
			Name:         recipeDB.name,
			Ingredients:  recipeIngredients,
			CreatedAt:    recipeDB.createdAt,
			LastModified: recipeDB.lastModified,
		})
	}
	return recipes, nil
}

func (r *recipeRepository) SaveRecipe(ctx context.Context, recipe *model.Recipe) error {
	return r.db.WithTx(ctx, func(tx database.TX) error {
		result, err := tx.ExecContext(ctx, "INSERT INTO recipe (name, created_at, last_modified) VALUES (?, ?, ?)", recipe.Name, recipe.CreatedAt, recipe.LastModified)
		if err != nil {
			return err
		}

		recipeID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for _, recipeIngredient := range recipe.Ingredients {
			_, err = tx.ExecContext(ctx, "INSERT INTO recipe_ingredient (recipe_id, ingredient_id, units) VALUES (?, ?, ?)", recipeID, recipeIngredient.Ingredient.ID, recipeIngredient.Units)
			if err != nil {
				return err
			}
		}
		recipe.ID = recipeID
		return nil
	})
}

func (r *recipeRepository) SaveRecipeSales(ctx context.Context, recipeSales *model.RecipeSales) error {
	return r.db.WithTx(ctx, func(tx database.TX) error {

		recipe, err := getRecipe(ctx, tx, recipeSales.RecipeID)

		if err != nil {
			return err
		}

		for _, recipeIngredient := range recipe.Ingredients {
			ingredientUsedUnits := recipeIngredient.Units * recipeSales.Units
			res, err := tx.ExecContext(ctx, "UPDATE ingredient SET units_in_stock = units_in_stock - ?, last_modified = ? WHERE id = ?", ingredientUsedUnits, recipeSales.CreatedAt, recipeIngredient.Ingredient.ID)
			if err != nil {
				return err
			}
			if rows, err := res.RowsAffected(); err != nil {
				return err
			} else if rows == 0 {
				return ErrNotFound
			}
		}

		result, err := tx.ExecContext(ctx, "INSERT INTO sold_recipes_history (recipe_id, units, created_at) VALUES (?, ?, ?)", recipeSales.RecipeID, recipeSales.Units, recipeSales.CreatedAt)

		if err != nil {
			r.logger.Error(err, "error saving ingredient stock")
			return err
		}

		recipeSalesID, err := result.LastInsertId()

		if err != nil {
			return err
		}

		recipeSales.ID = recipeSalesID
		return nil
	})
}
